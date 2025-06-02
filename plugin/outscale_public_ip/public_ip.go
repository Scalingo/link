package outscalepublicip

import (
	"context"
	"fmt"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/sirupsen/logrus"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/services/outscale"
)

type Plugin struct {
	oscClient outscale.PublicIPClient

	refreshEvery time.Duration

	// Public IP Configuration
	publicIPID string // ID of the public IP to move
	nicID      string // ID of the NIC to move the public IP to

	// Internal configuration
	linkPublicIPID  string // ID of the LinK from the public IP to the NIC
	lastRefreshedAt time.Time
}

func (p *Plugin) Activate(ctx context.Context) error {
	ctx, log := logger.WithStructToCtx(ctx, "plugin", p)

	log.Info("Linking public IP to NIC")
	resp, err := p.oscClient.LinkPublicIP(ctx, osc.LinkPublicIpRequest{
		PublicIpId:  &p.publicIPID,
		NicId:       &p.nicID,
		AllowRelink: osc.PtrBool(true),
	})

	if err != nil {
		return errors.Wrap(ctx, err, "link public IP")
	}

	p.linkPublicIPID = resp.GetLinkPublicIpId()
	p.lastRefreshedAt = time.Now()

	return nil
}

func (p *Plugin) Deactivate(ctx context.Context) error {
	ctx, log := logger.WithStructToCtx(ctx, "plugin", p)
	if p.linkPublicIPID == "" {
		log.Info("Public IP was not linked to a NIC, skipping unlink")
		return nil
	}

	log.WithField("link_public_ip_id", p.linkPublicIPID).Info("Unlinking public IP from NIC")

	_, err := p.oscClient.UnlinkPublicIP(ctx, osc.UnlinkPublicIpRequest{
		LinkPublicIpId: &p.linkPublicIPID,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "unlink public IP")
	}

	p.linkPublicIPID = ""

	return nil
}

func (p *Plugin) Ensure(ctx context.Context) error {
	ctx, log := logger.WithStructToCtx(ctx, "plugin", p)

	if p.lastRefreshedAt.Add(p.refreshEvery).After(time.Now()) {
		log.Debug("Already refreshed recently, skipping")
		return nil
	}

	publicIP, err := p.oscClient.ReadPublicIP(ctx, p.publicIPID)
	if err != nil {
		return errors.Wrap(ctx, err, "read public IP")
	}

	// If the public IP is not linked to the NIC, we need to link it
	if publicIP.GetNicId() != p.nicID {
		log.WithField("nic_id", publicIP.GetNicId()).Info("Public IP is not linked to the NIC, linking it")
		err := p.Activate(ctx)
		if err != nil {
			return errors.Wrap(ctx, err, "link public IP")
		}
		return nil
	}

	// The public IP is linked to the NIC, but the link ID is different
	if publicIP.GetLinkPublicIpId() != p.linkPublicIPID {
		// Update the link ID
		p.linkPublicIPID = publicIP.GetLinkPublicIpId()
	}

	p.lastRefreshedAt = time.Now()

	return nil
}

func (p *Plugin) ElectionKey(_ context.Context) string {
	return fmt.Sprintf("%s/%s", Name, p.publicIPID)
}

func (p *Plugin) LogFields() logrus.Fields {
	return logrus.Fields{
		"name":         "outscale_public_ip",
		"public_ip_id": p.publicIPID,
		"nic_id":       p.nicID,
	}
}
