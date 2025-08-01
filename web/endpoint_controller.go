package web

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/endpoint"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/scheduler"
)

var vipRegex = regexp.MustCompile(`^vip-[a-zA-Z0-9-]{36}$`)

type EndpointController struct {
	scheduler        scheduler.Scheduler
	storage          models.Storage
	endpointCreator  endpoint.Creator
	encryptedStorage models.EncryptedStorage
}

func NewEndpointController(scheduler scheduler.Scheduler, storage models.Storage, endpointCreator endpoint.Creator, encryptedStorage models.EncryptedStorage) EndpointController {
	return EndpointController{
		scheduler:        scheduler,
		storage:          storage,
		endpointCreator:  endpointCreator,
		encryptedStorage: encryptedStorage,
	}
}

func (c EndpointController) List(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	ips := c.scheduler.ConfiguredEndpoints(ctx)

	res := api.EndpointListResponse{
		Endpoints: ips.ToAPIType(),
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return errors.Wrap(ctx, err, "encode endpoints")
	}
	return nil
}

func (c EndpointController) Get(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	endpointId := params["id"]
	if !isEndpointIDValid(endpointId) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(ctx, "Invalid endpoint ID")
	}

	endpoint := c.scheduler.GetEndpoint(r.Context(), endpointId)
	if endpoint == nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New(ctx, "Endpoint not found")
	}

	err := json.NewEncoder(w).Encode(api.EndpointGetResponse{
		Endpoint: endpoint.ToAPIType(),
	})
	if err != nil {
		return errors.Wrap(ctx, err, "encode endpoint")
	}
	return nil
}

type CreateEndpointParams struct {
	HealthCheckInterval int               `json:"healthcheck_interval"`
	Checks              []api.HealthCheck `json:"checks"`
	Plugin              string            `json:"plugin"`
	PluginConfig        json.RawMessage   `json:"plugin_config"`
}

func (c EndpointController) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	var params CreateEndpointParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return errors.Wrap(ctx, err, "decode endpoint")
	}

	endpoint, err := c.endpointCreator.CreateEndpoint(ctx, endpoint.CreateEndpointParams{
		HealthCheckInterval: params.HealthCheckInterval,
		Checks:              params.Checks,
		Plugin:              params.Plugin,
		PluginConfig:        params.PluginConfig,
	})
	if err != nil {
		if errors.Is(err, scheduler.ErrEndpointAlreadyAssigned) {
			w.WriteHeader(http.StatusConflict)
			return err
		}
		return errors.Wrap(ctx, err, "create endpoint")
	}

	endpointWithStatus := c.scheduler.GetEndpoint(ctx, endpoint.ID)
	if endpointWithStatus == nil {
		return errors.New(ctx, "endpoint not found after creation")
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(api.EndpointGetResponse{
		Endpoint: endpointWithStatus.ToAPIType(),
	})
	if err != nil {
		return errors.Wrap(ctx, err, "encode endpoint")
	}

	return nil
}

func (c EndpointController) Update(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	endpointId := params["id"]
	if !isEndpointIDValid(endpointId) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(ctx, "Invalid endpoint ID")
	}

	endpoint := c.scheduler.GetEndpoint(ctx, endpointId)
	if endpoint == nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New(ctx, "Endpoint not found")
	}

	var patchParams api.UpdateEndpointParams
	err := json.NewDecoder(r.Body).Decode(&patchParams)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid JSON")
	}

	checks := models.HealthChecksFromAPIType(patchParams.HealthChecks)
	validationErr := checks.Validate(ctx)
	if validationErr != nil {
		return errors.Wrap(ctx, validationErr, "validate health checks")
	}
	endpoint.Checks = checks
	err = c.scheduler.UpdateEndpoint(ctx, endpoint.Endpoint)
	if err != nil {
		return errors.Wrap(ctx, err, "update endpoint")
	}

	err = json.NewEncoder(w).Encode(endpoint.ToAPIType())
	if err != nil {
		return errors.Wrap(ctx, err, "encode endpoint")
	}

	return nil
}

func (c EndpointController) Delete(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	id := params["id"]

	if !isEndpointIDValid(id) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(ctx, "Invalid endpoint ID")
	}

	err := c.storage.RemoveEndpoint(ctx, id)
	if err != nil {
		return errors.Wrap(ctx, err, "remove endpoint from storage")
	}

	err = c.encryptedStorage.Cleanup(ctx, id)
	if err != nil {
		return errors.Wrap(ctx, err, "cleanup encrypted storage for endpoint")
	}

	err = c.scheduler.Stop(ctx, id)
	if err != nil {
		return errors.Wrap(ctx, err, "stop endpoint manager")
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c EndpointController) Failover(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	id := params["id"]

	if !isEndpointIDValid(id) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(ctx, "Invalid endpoint ID")
	}

	err := c.scheduler.Failover(ctx, id)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to failover endpoint")
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c EndpointController) GetHosts(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx := r.Context()

	id := params["id"]
	if !isEndpointIDValid(id) {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(ctx, "Invalid endpoint ID")
	}

	endpoint := c.scheduler.GetEndpoint(ctx, id)
	if endpoint == nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New(ctx, "Endpoint not found")
	}

	hosts, err := c.storage.GetEndpointHosts(ctx, endpoint.ElectionKey)
	if err != nil {
		return errors.Wrap(ctx, err, "get endpoint hosts")
	}

	hostsRes := make([]api.Host, 0, len(hosts))
	for _, host := range hosts {
		hostsRes = append(hostsRes, api.Host{
			Hostname: host,
		})
	}

	response := api.GetEndpointHostsResponse{
		Hosts: hostsRes,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return errors.Wrap(ctx, err, "encode hosts")
	}

	return nil
}

func isEndpointIDValid(id string) bool {
	if id == "" {
		return false
	}

	if !vipRegex.MatchString(id) {
		return false
	}

	return true
}
