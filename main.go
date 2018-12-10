package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/logger/plugins/rollbarplugin"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/scheduler"
	"github.com/Scalingo/link/web"
	"github.com/sirupsen/logrus"
)

var Version = "dev"

func main() {
	config, err := config.Build()
	if err != nil {
		panic(err)
	}

	rollbarplugin.Register()
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	etcd, err := etcd.ClientFromEnv()
	if err != nil {
		panic(err)
	}

	storage := models.NewEtcdStorage(config)
	scheduler := scheduler.NewIPScheduler(config, etcd, storage)

	ips, err := storage.GetIPs(ctx)
	if err != nil {
		panic(err)
	}

	if len(ips) > 0 {
		log.Info("Restarting IP schedulers...")
		for _, ip := range ips {
			log := log.WithFields(logrus.Fields{
				"id": ip.ID,
				"ip": ip.IP,
			})
			log.Info("Starting")
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
	r.HandleFunc("/ips/{id}/lock", ipController.TryGetLock).Methods("POST")
	r.HandleFunc("/version", versionController.Version).Methods("GET")

	log.Infof("Listening on %v", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), r)
	if err != nil {
		panic(err)
	}
}
