package webhook

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cryptoutils "github.com/Scalingo/go-utils/crypto"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
)

type Plugin struct {
	endpoint   models.Endpoint
	cfg        PluginConfig
	httpClient *http.Client

	refreshEvery    time.Duration
	lastRefreshedAt time.Time
}

func (p *Plugin) Activate(ctx context.Context) error {
	log := logger.Get(ctx)
	payload, err := p.buildPayload(api.Activated)
	if err != nil {
		return errors.Wrap(ctx, err, "marshal webhook payload")
	}

	log.Info("Sending activation webhook")
	err = p.sendWebhook(ctx, payload)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook")
	}

	p.lastRefreshedAt = time.Now()
	log.Info("Activation webhook sent successfully")

	return nil
}

func (p *Plugin) Deactivate(ctx context.Context) error {
	log := logger.Get(ctx)
	payload, err := p.buildPayload(api.Standby)
	if err != nil {
		return errors.Wrap(ctx, err, "marshal webhook payload")
	}

	log.Info("Sending deactivation webhook")
	err = p.sendWebhook(ctx, payload)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook")
	}
	log.Info("Deactivation webhook sent successfully")

	return nil
}

func (p *Plugin) Ensure(ctx context.Context) error {
	log := logger.Get(ctx)
	if p.lastRefreshedAt.Add(p.refreshEvery).After(time.Now()) {
		log.Debug("No need to refresh webhook yet")
		return nil
	}

	err := p.Activate(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "activate webhook")
	}

	return nil
}

func (p *Plugin) ElectionKey(_ context.Context) string {
	return fmt.Sprintf("%s/%s", Name, p.cfg.ResourceID)
}

func (p *Plugin) buildPayload(status string) ([]byte, error) {
	return json.Marshal(api.WebhookPluginStatusChangePayload{
		EndpointID: p.endpoint.ID,
		ResourceID: p.cfg.ResourceID,
		Plugin:     p.endpoint.Plugin,
		Status:     status,
	})
}

func (p *Plugin) sendWebhook(ctx context.Context, payload []byte) error {
	log := logger.Get(ctx)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(ctx, err, "create webhook request")
	}

	req.Header.Set("Content-Type", "application/json")
	for name, value := range p.cfg.Headers {
		req.Header.Set(name, value)
	}

	if p.cfg.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		signature := cryptoutils.HMAC256([]byte(p.cfg.Secret), []byte(timestamp+"."+string(payload)))

		req.Header.Set(api.HeaderWebhookTimestamp, timestamp)
		req.Header.Set(api.HeaderWebhookSignature, hex.EncodeToString(signature))
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook request")
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.WithError(err).Error("Fail to close webhook response body")
		}
	}()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return errors.Newf(ctx, "webhook returned non-success status code: %d", res.StatusCode)
	}

	return nil
}
