package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	scalingoerrors "github.com/Scalingo/go-utils/errors"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/scheduler"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type ipController struct {
	scheduler scheduler.Scheduler
}

func NewIPController(scheduler scheduler.Scheduler) ipController {
	return ipController{
		scheduler: scheduler,
	}
}

func (c ipController) List(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
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

// Patch updates the healthchecks configured on the given IP. It actually **replaces** the currently configured healthchecks.
func (c ipController) Patch(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	ip := c.scheduler.GetIP(ctx, params["id"])
	if ip == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "not found"}`))
		return nil
	}
	log = log.WithFields(ip.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)
	log.Info("Updating a LinK IP healthchecks")

	var patchParams api.UpdateIPParams
	err := json.NewDecoder(r.Body).Decode(&patchParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid json"}`))
		return nil
	}

	err = checkIPHealthchecks(ctx, patchParams.Healthchecks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return nil
	}

	ip.Checks = patchParams.Healthchecks
	err = c.scheduler.UpdateIP(ctx, ip.IP)
	if err != nil {
		if scalingoerrors.RootCause(err) == scheduler.ErrManagerNotFound {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
			return nil
		}
		return errors.Wrapf(err, "fail to update the LinK IP '%s'", ip.ID)
	}

	err = json.NewEncoder(w).Encode(map[string]api.IP{
		"ip": *ip,
	})
	if err != nil {
		log.WithError(err).Error("fail to encode the IP after patching it")
		return nil
	}

	return nil
}

func (c ipController) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
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
	log = log.WithFields(ip.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)
	log.Info("Creating a new LinK IP")

	_, err = netlink.ParseAddr(ip.IP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg": "invalid IP"}`))
		return nil
	}
	err = checkIPHealthchecks(ctx, ip.Checks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return nil
	}

	ctx = logger.ToCtx(context.Background(), log)
	ip, err = c.scheduler.Start(ctx, ip)
	if err != nil {
		if errors.Cause(err) == scheduler.ErrNotStopping {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg": "IP already assigned"}`))
			return nil
		}
		return errors.Wrap(err, "fail to start IP manager")
	}

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
		return errors.Wrap(err, "fail to stop IP manager")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c ipController) TryGetLock(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	id := params["id"]

	w.Header().Set("Content-Type", "application/json")

	found := c.scheduler.TryGetLock(ctx, id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func checkIPHealthchecks(ctx context.Context, healthchecks []models.Healthcheck) error {
	for _, check := range healthchecks {
		if check.Port <= 0 {
			return errors.New("health check port cannot be negative")
		}
		if check.Port > 65535 {
			return errors.New("health check port cannot be greater than 65535")
		}
	}

	return nil
}
