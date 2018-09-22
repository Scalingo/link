package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/scheduler"
	"github.com/Scalingo/link/web"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.Build()
	if err != nil {
		panic(err)
	}

	log := logger.Default()
	log.SetLevel(logrus.InfoLevel)
	ctx := logger.ToCtx(context.Background(), log)

	etcd, err := etcd.ClientFromEnv()
	if err != nil {
		panic(err)
	}

	storage := models.NewETCDStorage(config)
	scheduler := scheduler.NewIPScheduler(config, etcd)

	ips, err := storage.GetIPs(ctx)
	if err != nil {
		panic(err)
	}

	if len(ips) > 0 {
		log.Info("Restarting IP schedulers...")
		for _, ip := range ips {
			log.WithFields(logrus.Fields{
				"id": ip.ID,
				"ip": ip.IP,
			}).Info("Starting")
			err := scheduler.Start(ctx, ip.ID, ip.IP)
			if err != nil {
				panic(err)
			}
		}
	}

	controller := web.NewIPController(storage, scheduler)
	r := handlers.NewRouter(log)
	r.Use(handlers.ErrorMiddleware)
	r.HandleFunc("/ips", controller.List).Methods("GET")
	r.HandleFunc("/ips", controller.Create).Methods("POST")
	r.HandleFunc("/ips/{id}", controller.Destroy).Methods("DELETE")

	log.Infof("Listening on %v", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), r)
	if err != nil {
		panic(err)
	}
}
