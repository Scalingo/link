package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/logger/plugins/rollbarplugin"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/locker"
	"github.com/Scalingo/link/v2/migrations"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/scheduler"
	"github.com/Scalingo/link/v2/web"
)

// Version is the current LinK version. During release build it will be overwritten by the compiler
var Version = "dev"

func main() {
	rollbarplugin.Register()
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	config, err := config.Build()
	if err != nil {
		log.WithError(err).Error("Fail to init config")
		panic(err)
	}

	etcd, err := etcd.ClientFromEnv()
	if err != nil {
		log.WithError(err).Error("Fail to get etcd client")
		panic(err)
	}

	storage := models.NewEtcdStorage(config)
	leaseManager := locker.NewEtcdLeaseManager(ctx, config, storage, etcd)

	// We need to check if it is needed to migrate data from v0 to v1 before the lease manager is started so that we are sure that the host data version is still the one from the previous LinK execution.
	migrationV0toV1 := migrations.NewV0toV1Migration(config.Hostname, leaseManager, storage)
	needsMigrationV0toV1, err := migrationV0toV1.NeedsMigration(ctx)
	if err != nil {
		panic(err)
	}

	err = leaseManager.Start(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to start lease manager")
		panic(err)
	}

	// v0 to v1 migration needs to be executed after the lease manager is started as we generate the host lease ID in this method.
	if needsMigrationV0toV1 {
		go func(ctx context.Context) {
			err := migrationV0toV1.Migrate(ctx)
			if err != nil {
				log.WithError(err).Error("Fail to migrate data from v0 to v1")
				return
			}
		}(ctx)
	}

	scheduler := scheduler.NewIPScheduler(config, etcd, storage, leaseManager)

	ips, err := storage.GetEndpoints(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to list configured IPs")
		panic(err)
	}

	if len(ips) > 0 {
		log.Info("Restarting IP schedulers...")
		for _, ip := range ips {
			log := log.WithFields(logrus.Fields{
				"id": ip.ID,
				"ip": ip.IP,
			})
			log.Info("Starting an IP scheduler")
			_, err := scheduler.Start(logger.ToCtx(ctx, log), ip)
			if err != nil {
				panic(err)
			}
		}
	}

	ipController := web.NewIPController(scheduler)
	versionController := web.NewVersionController(Version)
	r := handlers.NewRouter(log)

	if config.User != "" || config.Password != "" {
		r.Use(handlers.AuthMiddleware(func(user, password string) bool {
			return user == config.User && password == config.Password
		}))
	}

	r.Use(handlers.ErrorMiddleware)
	r.HandleFunc("/ips", ipController.List).Methods("GET")
	r.HandleFunc("/ips", ipController.Create).Methods("POST")
	r.HandleFunc("/ips/{id}", ipController.Destroy).Methods("DELETE")
	r.HandleFunc("/ips/{id}", ipController.Get).Methods("GET")
	r.HandleFunc("/ips/{id}", ipController.Patch).Methods("PUT", "PATCH")
	r.HandleFunc("/ips/{id}/failover", ipController.Failover).Methods("POST")
	r.HandleFunc("/version", versionController.Version).Methods("GET")

	globalRouter := mux.NewRouter()

	if os.Getenv("PPROF_ENABLED") == "true" {
		pprofPrefix := "/debug/pprof"
		log.Info("Enabling pprof endpoints under " + pprofPrefix)

		pprofRouter := mux.NewRouter()
		pprofRouter.HandleFunc(pprofPrefix+"/", pprof.Index)
		pprofRouter.HandleFunc(pprofPrefix+"/profile", pprof.Profile)
		pprofRouter.HandleFunc(pprofPrefix+"/symbol", pprof.Symbol)
		pprofRouter.HandleFunc(pprofPrefix+"/cmdline", pprof.Cmdline)
		pprofRouter.HandleFunc(pprofPrefix+"/trace", pprof.Trace)
		pprofRouter.Handle(pprofPrefix+"/heap", pprof.Handler("heap"))
		pprofRouter.Handle(pprofPrefix+"/goroutine", pprof.Handler("goroutine"))
		pprofRouter.Handle(pprofPrefix+"/threadcreate", pprof.Handler("threadcreate"))
		pprofRouter.Handle(pprofPrefix+"/block", pprof.Handler("block"))

		globalRouter.Handle(pprofPrefix+"/{prop:.*}", pprofRouter)
	}
	globalRouter.Handle("/{any:.+}", r)

	log.Infof("Listening on %v", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), globalRouter)
	if err != nil {
		panic(err)
	}
}
