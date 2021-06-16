package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/ip"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/scheduler"
)

type ipController struct {
	scheduler scheduler.Scheduler
}

func NewIPController(scheduler scheduler.Scheduler) ipController {
	return ipController{
		scheduler: scheduler,
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
		log.WithError(err).Error("fail to encode IPs")
		return nil
	}
	return nil
}

func (c ipController) Get(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	ip := c.scheduler.GetIP(ctx, params["id"])
	if ip == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
		return nil
	}

	err := json.NewEncoder(w).Encode(map[string]api.IP{
		"ip": *ip,
	})

	if err != nil {
		log.WithError(err).Error("fail to encode IP")
		return nil
	}

	return nil
}

func (c ipController) Create(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)

	w.Header().Set("Content-Type", "application/json")
	var ip models.IP
	err := json.NewDecoder(r.Body).Decode(&ip)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "invalid json"}`))
		return nil
	}
	ip.ID = ""

	_, err = netlink.ParseAddr(ip.IP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "invalid IP"}`))
		return nil
	}
	for _, check := range ip.Checks {
		if check.Port <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg": "health check port cannot be negative"}`))
			return nil
		}
		if check.Port > 65535 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg": "health check port cannot be greater than 65535"}`))
			return nil
		}
	}

	ctx = logger.ToCtx(context.Background(), log)
	ip, err = c.scheduler.Start(ctx, ip)
	if err != nil {
		if err == scheduler.ErrIPAlreadyAssigned {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg": "IP already assigned"}`))
			return nil
		}
		return errors.Wrap(err, "fail to start IP manager")
	}
	log = log.WithFields(logrus.Fields{
		"id": ip.ID,
		"ip": ip.IP,
	})

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(ip)
	if err != nil {
		log.WithError(err).Error("fail to encode IP")
	}
	return nil
}

func (c ipController) Destroy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	id := params["id"]
	err := c.scheduler.Stop(ctx, id)
	if err != nil {
		if err == scheduler.ErrIPNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"msg": "IP not found"}`))
			return nil
		}
		return errors.Wrap(err, "fail to stop IP manager")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c ipController) Failover(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	id := params["id"]
	log := logger.Get(ctx).WithField("vip_id", id)
	err := c.scheduler.Failover(ctx, id)
	if err != nil {
		if err == scheduler.ErrIPNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"msg": "IP not found"}`))
			return nil
		}
		cause := errors.Cause(err)
		if cause == ip.ErrIsNotMaster || cause == ip.ErrNoOtherHosts {
			log.WithError(err).Info("Bad request: cannot failover")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"msg": cause.Error(),
			})
			return nil
		}
		return errors.Wrap(err, "fail to stop IP manager")
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
