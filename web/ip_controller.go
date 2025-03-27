package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/ip"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/scheduler"
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

	ips := c.scheduler.ConfiguredEndpoints(ctx)

	res := api.EndpointListResponse{
		Endpoints: ips.ToAPIType(),
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.WithError(err).Error("Fail to encode IPs")
		return nil
	}
	return nil
}

func (c ipController) Get(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	ip := c.scheduler.GetEndpoint(ctx, params["id"])
	if ip == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"resource": "IP", "error": "not found"}`))
		return nil
	}

	err := json.NewEncoder(w).Encode(api.EndpointGetResponse{
		Endpoint: ip.ToAPIType(),
	})

	if err != nil {
		log.WithError(err).Error("Fail to encode IP")
		return nil
	}

	return nil
}

func (c ipController) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)

	w.Header().Set("Content-Type", "application/json")
	var endpointParams api.AddEndpointParams
	err := json.NewDecoder(r.Body).Decode(&endpointParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid json"}`))
		return nil
	}
	_, err = netlink.ParseAddr(endpointParams.IP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid IP"}`))
		return nil
	}

	err = checkEndpointHealthChecks(ctx, endpointParams.Checks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return nil
	}

	endpoint := models.Endpoint{
		IP:                  endpointParams.IP,
		HealthCheckInterval: endpointParams.HealthCheckInterval,
		Checks:              models.HealthChecksFromAPIType(endpointParams.Checks),
	}

	endpoint.ID = ""
	log = log.WithFields(endpoint.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)
	log.Info("Creating a new LinK Endpoint")

	ctx = logger.ToCtx(context.Background(), log)
	endpoint, err = c.scheduler.Start(ctx, endpoint)
	if err != nil {
		if err == scheduler.ErrIPAlreadyAssigned {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "IP already assigned"}`))
			return nil
		}
		return errors.Wrap(err, "start endpoint manager")
	}
	log = log.WithFields(endpoint.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(endpoint.ToAPIType())
	if err != nil {
		log.WithError(err).Error("Fail to encode endpoint")
	}
	return nil
}

// Patch updates the health checks configured on the given IP. It actually **replaces** the currently configured healthchecks.
func (c ipController) Patch(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	endpoint := c.scheduler.GetEndpoint(ctx, params["id"])
	if endpoint == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"resource": "IP", "error": "not found"}`))
		return nil
	}
	log = log.WithFields(endpoint.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)
	log.Info("Updating an endpoint health checks")

	var patchParams api.UpdateEndpointParams
	err := json.NewDecoder(r.Body).Decode(&patchParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid json"}`))
		return nil
	}

	err = checkEndpointHealthChecks(ctx, patchParams.HealthChecks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return nil
	}

	endpoint.Checks = models.HealthChecksFromAPIType(patchParams.HealthChecks)
	err = c.scheduler.UpdateEndpoint(ctx, endpoint.Endpoint)
	if err != nil {
		return errors.Wrapf(err, "fail to update the LinK IP '%s'", endpoint.ID)
	}

	err = json.NewEncoder(w).Encode(endpoint.ToAPIType())
	if err != nil {
		log.WithError(err).Error("Fail to encode the IP after patching it")
		return nil
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
			w.Write([]byte(`{"resource": "IP", "error": "not found"}`))
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
			w.Write([]byte(`{"resource": "IP", "error": "not found"}`))
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

func checkEndpointHealthChecks(ctx context.Context, healthChecks []api.HealthCheck) error {
	for _, check := range healthChecks {
		if check.Port <= 0 {
			return errors.New("health check port cannot be negative")
		}
		if check.Port > 65535 {
			return errors.New("health check port cannot be greater than 65535")
		}
	}

	return nil
}
