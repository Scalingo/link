package utils

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/link/v2/api"
)

func GetClient(c *cli.Command) api.HTTPClient {
	var opts []api.ClientOpt
	if c.String("host") != "" {
		opts = append(opts, api.WithURL(c.String("host")))
	}

	if c.String("user") != "" {
		opts = append(opts, api.WithUser(c.String("user")))
	}

	if c.String("password") != "" {
		opts = append(opts, api.WithPassword(c.String("password")))
	}

	return api.NewHTTPClient(opts...)
}
