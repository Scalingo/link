package web

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/endpoint"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin/arp"
	"github.com/Scalingo/link/v3/scheduler"
)

// Legacy types. Those should be removed later in the v3 branch
type AddIPParams struct {
	HealthCheckInterval int               `json:"healthcheck_interval"`
	Checks              []api.HealthCheck `json:"checks"`
	IP                  string            `json:"ip"`
}

type ListIPsResponse struct {
	IPs []api.Endpoint `json:"ips"`
}

type GetIPResponse struct {
	IP api.Endpoint `json:"ip"`
}

type IPController struct {
	scheduler       scheduler.Scheduler
	storage         models.Storage
	endpointCreator endpoint.Creator
}

func NewIPController(scheduler scheduler.Scheduler, storage models.Storage, endpointCreator endpoint.Creator) IPController {
	return IPController{
		scheduler:       scheduler,
		storage:         storage,
		endpointCreator: endpointCreator,
	}
}

func (c IPController) List(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	ips := c.scheduler.ConfiguredEndpoints(ctx)

	res := ListIPsResponse{
		IPs: ips.ToAPIType(),
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return errors.Wrap(ctx, err, "encode endpoints")
	}
	return nil
}

func (c IPController) Get(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)
	w.Header().Set("Content-Type", "application/json")

	ip := c.scheduler.GetEndpoint(ctx, params["id"])
	if ip == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"resource": "IP", "error": "not found"}`))
		return nil
	}

	err := json.NewEncoder(w).Encode(GetIPResponse{
		IP: ip.ToAPIType(),
	})

	if err != nil {
		log.WithError(err).Error("Fail to encode IP")
		return nil
	}

	return nil
}

func (c IPController) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	log := logger.Get(ctx)

	w.Header().Set("Content-Type", "application/json")
	var endpointParams AddIPParams
	err := json.NewDecoder(r.Body).Decode(&endpointParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid json"}`))
		return nil
	}

	pluginConfig, _ := json.Marshal(arp.PluginConfig{
		IP: endpointParams.IP,
	})

	endpoint, err := c.endpointCreator.CreateEndpoint(ctx, endpoint.CreateEndpointParams{
		HealthCheckInterval: endpointParams.HealthCheckInterval,
		Checks:              endpointParams.Checks,
		Plugin:              "arp",
		PluginConfig:        pluginConfig,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "create endpoint")
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(endpoint.ToAPIType())
	if err != nil {
		log.WithError(err).Error("Fail to encode endpoint")
	}
	return nil
}
