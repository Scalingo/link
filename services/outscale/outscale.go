package outscale

import (
	"context"

	"github.com/outscale/osc-sdk-go/v2"
)

// The outscale package provides a generic interface for interacting with the Outscale API.
// This is useful for generating mocks and testing purposes.
// It also helps with abstracting some of the complexities of the Outscale API.

type APIClient struct {
	oscClient *osc.APIClient

	// Client Configuration
	accessKey string
	secretKey string
	region    string
}

func NewClient(accessKey, secretKey, region string) *APIClient {
	oscClient := osc.NewAPIClient(osc.NewConfiguration())

	return &APIClient{
		oscClient: oscClient,
		accessKey: accessKey,
		secretKey: secretKey,
		region:    region,
	}
}

func (c *APIClient) authenticatedContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, osc.ContextAWSv4, osc.AWSv4{
		AccessKey: c.accessKey,
		SecretKey: c.secretKey,
	})

	ctx = context.WithValue(ctx, osc.ContextServerIndex, 0)
	ctx = context.WithValue(ctx, osc.ContextServerVariables, map[string]string{
		"region": c.region,
	})
	return ctx
}
