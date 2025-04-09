package arp

import (
	"context"
	"encoding/json"

	"github.com/kelseyhightower/envconfig"
	"github.com/vishvananda/netlink"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network"
	"github.com/Scalingo/link/v2/plugin"
)

const Name = "arp"

type Config struct {
	// Number of gratuitous ARP (GARP) packets sent when the state becomes 'ACTIVATED'
	ARPGratuitousCount int    `envconfig:"ARP_GRATUITOUS_COUNT" default:"3"`
	Interface          string `envconfig:"INTERFACE"`
}

func Register(ctx context.Context, registry plugin.Registry) error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.Wrap(ctx, err, "parse environment")
	}

	i, err := network.NewNetworkInterfaceFromName(config.Interface)
	if err != nil {
		return errors.Wrap(ctx, err, "get network interface")
	}

	registry.Register(ctx, Name, Factory{
		config:       config,
		netInterface: i,
	})

	return nil
}

type Factory struct {
	config       Config
	netInterface network.Interface
}

type PluginConfig struct {
	IP string `json:"ip"`
}

func (f Factory) Create(ctx context.Context, endpoint models.Endpoint) (plugin.Plugin, error) {
	var cfg PluginConfig

	if endpoint.PluginConfig != nil {
		err := json.Unmarshal(endpoint.PluginConfig, &cfg)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "unmarshal plugin config")
		}
	} else { // Retro compatibility
		if endpoint.IP == "" {
			return nil, errors.New(ctx, "invalid plugin config: empty")
		}
		cfg.IP = endpoint.IP
	}

	return &Plugin{
		endpoint:     endpoint,
		ip:           cfg.IP,
		config:       f.config,
		garpCount:    0,
		netInterface: f.netInterface,
	}, nil
}

func (f Factory) Validate(_ context.Context, endpoint models.Endpoint) error {
	validation := errors.NewValidationErrorsBuilder()
	var cfg PluginConfig
	err := json.Unmarshal(endpoint.PluginConfig, &cfg)
	if err != nil {
		validation.Set("plugin_config", "invalid JSON: "+err.Error())
		return validation.Build()
	}

	if cfg.IP == "" {
		validation.Set("plugin_config.ip", "ip is required")
	} else {
		_, err = netlink.ParseAddr(cfg.IP)
		if err != nil {
			validation.Set("plugin_config.ip", "invalid IP address: "+err.Error())
		}
	}

	validationErr := validation.Build()
	if validationErr != nil {
		return validationErr
	}
	return nil
}
