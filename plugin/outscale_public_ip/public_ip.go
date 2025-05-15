package outscalepublicip

import (
	"context"
	"fmt"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/sirupsen/logrus"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

type Plugin struct {
	oscClient *osc.APIClient

	refreshEvery time.Duration

	// Client Configuration
	accessKey string
	secretKey string
	region    string

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
	resp, _, err := p.oscClient.PublicIpApi.LinkPublicIp(p.authenticatedContext(ctx)).LinkPublicIpRequest(osc.LinkPublicIpRequest{
		PublicIpId:  &p.publicIPID,
		NicId:       &p.nicID,
		AllowRelink: osc.PtrBool(true),
	}).Execute()
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

	_, _, err := p.oscClient.PublicIpApi.UnlinkPublicIp(p.authenticatedContext(ctx)).UnlinkPublicIpRequest(osc.UnlinkPublicIpRequest{
		LinkPublicIpId: &p.linkPublicIPID,
	}).Execute()
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

	resp, _, err := p.oscClient.PublicIpApi.ReadPublicIps(p.authenticatedContext(ctx)).ReadPublicIpsRequest(osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{
			PublicIpIds: &[]string{p.publicIPID},
		},
	}).Execute()
	if err != nil {
		return errors.Wrap(ctx, err, "read public IP")
	}

	// Check if the PublicIP response is sane
	if len(resp.GetPublicIps()) == 0 {
		return errors.New(ctx, "public IP not found")
	}

	if len(resp.GetPublicIps()) > 1 {
		return errors.New(ctx, "multiple public IPs found")
	}

	publicIP := resp.GetPublicIps()[0]

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

func (p *Plugin) authenticatedContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, osc.ContextAWSv4, osc.AWSv4{
		AccessKey: p.accessKey,
		SecretKey: p.secretKey,
	})

	ctx = context.WithValue(ctx, osc.ContextServerIndex, 0)
	ctx = context.WithValue(ctx, osc.ContextServerVariables, map[string]string{
		"region": p.region,
	})
	return ctx
}

func (p *Plugin) LogFields() logrus.Fields {
	return logrus.Fields{
		"name":         "outscale_public_ip",
		"public_ip_id": p.publicIPID,
		"nic_id":       p.nicID,
	}
}
