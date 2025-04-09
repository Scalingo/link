package arp

import (
	"context"
	"strings"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network"
)

type Plugin struct {
	config       Config
	netInterface network.Interface

	endpoint models.Endpoint
	ip       string

	garpCount int
}

func (p *Plugin) Activate(ctx context.Context) error {
	p.garpCount = 0
	err := p.netInterface.EnsureIP(p.ip)
	if err != nil {
		return errors.Wrap(ctx, err, "activate IP on network interface")
	}
	return nil
}

func (p *Plugin) Deactivate(ctx context.Context) error {
	p.garpCount = 0
	err := p.netInterface.RemoveIP(p.ip)
	if err != nil {
		return errors.Wrap(ctx, err, "disable IP on network interface")
	}
	return nil
}

func (p *Plugin) Ensure(ctx context.Context) error {
	ctx, log := logger.WithFieldToCtx(ctx, "plugin", "arp")

	if p.garpCount >= p.config.ARPGratuitousCount {
		log.Debug("All gratuitous ARP requests sent for this activation")
		return nil
	}

	log.Info("Send gratuitous ARP request")
	err := p.netInterface.EnsureIP(p.ip)
	if err != nil {
		return errors.Wrap(ctx, err, "send gratuitous ARP")
	}
	p.garpCount++

	return nil
}

func (p *Plugin) ElectionKey(_ context.Context) string {
	return strings.ReplaceAll(p.ip, "/", "_")
}
