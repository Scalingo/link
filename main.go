package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gorilla/mux"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/logger/plugins/rollbarplugin"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/endpoint"
	"github.com/Scalingo/link/v3/locker"
	"github.com/Scalingo/link/v3/migrations"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin"
	"github.com/Scalingo/link/v3/plugin/arp"
	outscalepublicip "github.com/Scalingo/link/v3/plugin/outscale_public_ip"
	"github.com/Scalingo/link/v3/scheduler"
	"github.com/Scalingo/link/v3/web"
)

// Version is the current LinK version. During release build it will be overwritten by the compiler
var Version = "dev"

func main() {
	rollbarplugin.Register()
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	config, err := config.Build(ctx)
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
	encryptedStorage, err := models.NewEncryptedStorage(ctx, config, storage)
	if err != nil {
		log.WithError(err).Error("Fail to init encrypted storage")
		panic(err)
	}

	pluginRegistry := plugin.NewRegistry()
	err = initPlugins(ctx, pluginRegistry, encryptedStorage)
	if err != nil {
		log.WithError(err).Error("Fail to init plugins")
		panic(err)
	}

	migrationRunner, err := migrations.NewMigrationRunner(ctx, config, storage, leaseManager)
	if err != nil {
		log.WithError(err).Error("Fail to init migration runner")
		panic(err)
	}

	// We run the migration in a goroutine. Because the migrations can take a long time and locks might expires.
	// This could cause unwanted failover.
	// All major versions of LinK should be compatible with the previous and next major data version.
	go func(ctx context.Context) {
		err := migrationRunner.Run(ctx)
		if err != nil {
			log.WithError(err).Error("Fail to run migrations")
			return
		}
	}(ctx)

	err = leaseManager.Start(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to start lease manager")
		panic(err)
	}

	scheduler := scheduler.NewEndpointScheduler(config, etcd, storage, leaseManager, pluginRegistry)

	endpoints, err := storage.GetEndpoints(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to list configured endpoints")
		panic(err)
	}

	if len(endpoints) > 0 {
		log.Info("Restarting endpoint schedulers...")
		for _, endpoint := range endpoints {
			ctx, log := logger.WithStructToCtx(ctx, "endpoint", endpoint)
			log.Info("Starting an endpoint scheduler")
			_, err := scheduler.Start(ctx, endpoint)
			if err != nil {
				panic(err)
			}
		}
	}

	endpointCreator := endpoint.NewCreator(storage, scheduler, pluginRegistry)
	ipController := web.NewIPController(scheduler, storage, endpointCreator)
	endpointController := web.NewEndpointController(scheduler, storage, endpointCreator)
	versionController := web.NewVersionController(Version)
	r := handlers.NewRouter(log)
	r.Use(handlers.ErrorMiddleware)

	if config.User != "" || config.Password != "" {
		r.Use(handlers.AuthMiddleware(func(user, password string) bool {
			return user == config.User && password == config.Password
		}))
	}

	r.Use(handlers.ErrorMiddleware)

	// Retro compatibility with v2 API.
	// This will be removed in a future version.
	r.HandleFunc("/ips", ipController.List).Methods("GET")
	r.HandleFunc("/ips", ipController.Create).Methods("POST")
	r.HandleFunc("/ips/{id}", endpointController.Delete).Methods("DELETE")
	r.HandleFunc("/ips/{id}", ipController.Get).Methods("GET")
	r.HandleFunc("/ips/{id}", endpointController.Update).Methods("PUT", "PATCH")
	r.HandleFunc("/ips/{id}/failover", endpointController.Failover).Methods("POST")

	r.HandleFunc("/endpoints", endpointController.List).Methods("GET")
	r.HandleFunc("/endpoints", endpointController.Create).Methods("POST")
	r.HandleFunc("/endpoints/{id}", endpointController.Delete).Methods("DELETE")
	r.HandleFunc("/endpoints/{id}", endpointController.Get).Methods("GET")
	r.HandleFunc("/endpoints/{id}", endpointController.Update).Methods("PUT", "PATCH")
	r.HandleFunc("/endpoints/{id}/failover", endpointController.Failover).Methods("POST")
	r.HandleFunc("/endpoints/{id}/hosts", endpointController.GetHosts).Methods("GET")

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

func initPlugins(ctx context.Context, registry plugin.Registry, encryptedStorage models.EncryptedStorage) error {
	err := arp.Register(ctx, registry)
	if err != nil {
		return errors.Wrap(ctx, err, "register arp plugin")
	}

	err = outscalepublicip.Register(ctx, registry, encryptedStorage)
	if err != nil {
		return errors.Wrap(ctx, err, "register outscale public ip plugin")
	}

	return nil
}
