package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/scheduler"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type IPController struct {
	storage   models.Storage
	scheduler scheduler.Scheduler
}

func NewIPController(storage models.Storage, scheduler scheduler.Scheduler) IPController {
	return IPController{
		storage:   storage,
		scheduler: scheduler,
	}
}

func (c IPController) List(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	ctx := r.Context()
	ips, err := c.storage.GetIPs(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get storage")
	}

	err = json.NewEncoder(w).Encode(map[string][]models.IP{
		"ips": ips,
	})
	if err != nil {
		return errors.Wrap(err, "fail to encode ips")
	}
	return nil
}

func (c IPController) Create(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	var newIP models.IP
	err := json.NewDecoder(r.Body).Decode(&newIP)
	if err != nil {
		return errors.Wrap(err, "invalid json")
	}

	_, err = netlink.ParseAddr(newIP.IP)
	if err != nil {
		return errors.Wrap(err, "Invalid IP")
	}

	newIP, err = c.storage.AddIP(ctx, newIP)
	if err != nil {
		return errors.Wrap(err, "fail to save IP")
	}

	ctx = logger.ToCtx(context.Background(), log)

	err = c.scheduler.Start(ctx, newIP.ID, newIP.IP)
	if err != nil {
		return errors.Wrap(err, "fail to start IP manager")
	}

	err = json.NewEncoder(w).Encode(newIP)
	if err != nil {
		return errors.Wrap(err, "fail to encode IP")
	}
	return nil
}
func (c IPController) Destroy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

	return nil
}
