package outscale

import (
	"context"

	"github.com/outscale/osc-sdk-go/v2"

	"github.com/Scalingo/go-utils/errors/v2"
)

var _ PublicIPClient = (*APIClient)(nil)

type PublicIPClient interface {
	LinkPublicIP(ctx context.Context, params osc.LinkPublicIpRequest) (osc.LinkPublicIpResponse, error)
	UnlinkPublicIP(ctx context.Context, params osc.UnlinkPublicIpRequest) (osc.UnlinkPublicIpResponse, error)
	ReadPublicIP(ctx context.Context, publicIPID string) (osc.PublicIp, error)
}

func (c *APIClient) LinkPublicIP(ctx context.Context, params osc.LinkPublicIpRequest) (osc.LinkPublicIpResponse, error) {
	authCtx := c.authenticatedContext(ctx)
	resp, _, err := c.oscClient.PublicIpApi.LinkPublicIp(authCtx).LinkPublicIpRequest(params).Execute()
	if err != nil {
		return osc.LinkPublicIpResponse{}, errors.Wrap(ctx, err, "link public IP")
	}
	return resp, nil
}

func (c *APIClient) UnlinkPublicIP(ctx context.Context, params osc.UnlinkPublicIpRequest) (osc.UnlinkPublicIpResponse, error) {
	authCtx := c.authenticatedContext(ctx)
	resp, _, err := c.oscClient.PublicIpApi.UnlinkPublicIp(authCtx).UnlinkPublicIpRequest(params).Execute()
	if err != nil {
		return osc.UnlinkPublicIpResponse{}, errors.Wrap(ctx, err, "unlink public IP")
	}
	return resp, nil
}

func (c *APIClient) ReadPublicIP(ctx context.Context, publicIPID string) (osc.PublicIp, error) {
	authCtx := c.authenticatedContext(ctx)
	req := osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{
			PublicIpIds: &[]string{publicIPID},
		},
	}
	resp, _, err := c.oscClient.PublicIpApi.ReadPublicIps(authCtx).ReadPublicIpsRequest(req).Execute()
	if err != nil {
		return osc.PublicIp{}, errors.Wrap(ctx, err, "read public IP")
	}
	if resp.PublicIps == nil {
		return osc.PublicIp{}, errors.New(ctx, "public IP not found")
	}

	if len(*resp.PublicIps) != 1 {
		return osc.PublicIp{}, errors.Newf(ctx, "invalid number of public IPs returned: %d", len(*resp.PublicIps))
	}

	return (*resp.PublicIps)[0], nil
}
