package scheduler

import (
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/models"
)

type EndpointWithStatus struct {
	models.Endpoint

	Status string
}

func (e EndpointWithStatus) ToAPIType() api.Endpoint {
	res := e.Endpoint.ToAPIType()
	res.Status = e.Status
	return res
}

type EndpointsWithStatus []EndpointWithStatus

func (e EndpointsWithStatus) ToAPIType() []api.Endpoint {
	res := make([]api.Endpoint, len(e))
	for i, ep := range e {
		res[i] = ep.ToAPIType()
	}
	return res
}
