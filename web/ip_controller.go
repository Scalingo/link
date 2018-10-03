package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/network"
	"github.com/Scalingo/link/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type ipController struct {
	storage      models.Storage
	scheduler    scheduler.Scheduler
	netInterface network.NetworkInterface
}

func NewIPController(storage models.Storage, scheduler scheduler.Scheduler, netInterface network.NetworkInterface) ipController {
	return ipController{
		storage:      storage,
		scheduler:    scheduler,
		netInterface: netInterface,
	}
}

func (c ipController) List(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	ips := c.scheduler.ConfiguredIPs(ctx)

	err := json.NewEncoder(w).Encode(map[string][]api.IP{
		"ips": ips,
	})
	if err != nil {
		log.WithError(err).Error(err, "fail to encode IPs")
	}
	return nil
}

func (c ipController) Create(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)

	w.Header().Set("Content-Type", "application/json")
	var newIP models.IP
	err := json.NewDecoder(r.Body).Decode(&newIP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "invalid json"}`))
		return nil
	}

	_, err = netlink.ParseAddr(newIP.IP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "invalid IP"}`))
		return nil
	}

	has, err := c.netInterface.HasIP(newIP.IP)
	if err != nil {
		return errors.Wrap(err, "fail to check if IP is already assigned")
	}

	if has {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "IP already assigned"}`))
		return nil
	}

	newIP, err = c.storage.AddIP(ctx, newIP)
	if err != nil {
		return errors.Wrap(err, "fail to save IP")
	}

	log = log.WithFields(logrus.Fields{
		"id": newIP.ID,
		"ip": newIP.IP,
	})

	ctx = logger.ToCtx(context.Background(), log)

	err = c.scheduler.Start(ctx, newIP)
	if err != nil {
		return errors.Wrap(err, "fail to start IP manager")
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newIP)
	if err != nil {
		log.WithError(err).Error("fail to encode IP")
	}
	return nil
}
func (c ipController) Destroy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	id := params["id"]
	err := c.storage.RemoveIP(ctx, id)
	if err != nil {
		return errors.Wrap(err, "fail to delete IP")
	}

	err = c.scheduler.Stop(ctx, id)
	if err != nil {
		return errors.Wrap(err, "fail to stop IP manager")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
