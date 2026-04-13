package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v3/models"
)

type Plugin struct {
	endpoint   models.Endpoint
	cfg        PluginConfig
	httpClient *http.Client
}

type statusChangePayload struct {
	EndpointID string    `json:"endpoint_id"`
	Plugin     string    `json:"plugin"`
	Status     string    `json:"status"`
	ChangedAt  time.Time `json:"changed_at"`
}

func (p *Plugin) Activate(ctx context.Context) error {
	payload, err := p.buildPayload(api.Activated)
	if err != nil {
		return errors.Wrap(ctx, err, "marshal webhook payload")
	}

	err = p.sendWebhook(ctx, payload)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook")
	}

	return nil
}

func (p *Plugin) Deactivate(ctx context.Context) error {
	payload, err := p.buildPayload(api.Standby)
	if err != nil {
		return errors.Wrap(ctx, err, "marshal webhook payload")
	}

	err = p.sendWebhook(ctx, payload)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook")
	}

	return nil
}

func (p *Plugin) Ensure(ctx context.Context) error {
	return p.Activate(ctx)
}

func (p *Plugin) ElectionKey(_ context.Context) string {
	return p.endpoint.ID
}

func (p *Plugin) buildPayload(status string) ([]byte, error) {
	return json.Marshal(statusChangePayload{
		EndpointID: p.endpoint.ID,
		Plugin:     p.endpoint.Plugin,
		Status:     status,
		ChangedAt:  time.Now().UTC(),
	})
}

func (p *Plugin) sendWebhook(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(ctx, err, "create webhook request")
	}

	req.Header.Set("Content-Type", "application/json")
	for name, value := range p.cfg.Headers {
		req.Header.Set(name, value)
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(ctx, err, "send webhook request")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return errors.Newf(ctx, "webhook returned non-success status code: %d", res.StatusCode)
	}

	return nil
}
