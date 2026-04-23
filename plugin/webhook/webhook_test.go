package webhook

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	cryptoutils "github.com/Scalingo/go-utils/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
)

func TestPluginOnStatusChange(t *testing.T) {
	t.Run("activate sends payload and configured headers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "token-123", r.Header.Get("Authorization"))
			assert.Equal(t, "link", r.Header.Get("X-App"))

			payloadBytes, err := io.ReadAll(r.Body)
			if !assert.NoError(t, err) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var body api.WebhookPluginStatusChangePayload
			err = json.Unmarshal(payloadBytes, &body)
			if !assert.NoError(t, err) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			assert.Equal(t, "vip-1", body.EndpointID)
			assert.Equal(t, "resource-activate", body.ResourceID)
			assert.Equal(t, Name, body.Plugin)
			assert.Equal(t, api.Activated, body.Status)
			timestamp := r.Header.Get(api.HeaderWebhookTimestamp)
			if !assert.NotEmpty(t, timestamp) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			assert.Equal(
				t,
				hex.EncodeToString(cryptoutils.HMAC256([]byte("webhook-secret"), []byte(timestamp+"."+string(payloadBytes)))),
				r.Header.Get(api.HeaderWebhookSignature),
			)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint: models.Endpoint{ID: "vip-1", Plugin: Name},
			cfg: PluginConfig{
				URL:        server.URL,
				ResourceID: "resource-activate",
				Secret:     "webhook-secret",
				Headers: map[string]string{
					"Authorization": "token-123",
					"X-App":         "link",
				},
			},
			httpClient: server.Client(),
		}

		err := p.Activate(t.Context())
		require.NoError(t, err)
	})

	t.Run("deactivate returns error on non 2xx status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint:   models.Endpoint{ID: "vip-2", Plugin: Name},
			cfg:        PluginConfig{URL: server.URL, ResourceID: "resource-default"},
			httpClient: server.Client(),
		}

		err := p.Deactivate(t.Context())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "webhook returned non-success status code")
	})
}

func TestPluginEnsure(t *testing.T) {
	t.Run("already refreshed recently", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint:        models.Endpoint{ID: "vip-ensure-1", Plugin: Name},
			cfg:             PluginConfig{URL: server.URL, ResourceID: "resource-default"},
			httpClient:      server.Client(),
			refreshEvery:    time.Minute,
			lastRefreshedAt: time.Now(),
		}

		err := p.Ensure(t.Context())
		require.NoError(t, err)
		assert.Equal(t, int32(0), calls.Load())
	})

	t.Run("send once in refresh window", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint:     models.Endpoint{ID: "vip-ensure-2", Plugin: Name},
			cfg:          PluginConfig{URL: server.URL, ResourceID: "resource-default"},
			httpClient:   server.Client(),
			refreshEvery: 30 * time.Minute,
		}

		err := p.Ensure(t.Context())
		require.NoError(t, err)
		err = p.Ensure(t.Context())
		require.NoError(t, err)

		assert.Equal(t, int32(1), calls.Load())
	})

	t.Run("refresh again after interval", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint:     models.Endpoint{ID: "vip-ensure-3", Plugin: Name},
			cfg:          PluginConfig{URL: server.URL, ResourceID: "resource-default"},
			httpClient:   server.Client(),
			refreshEvery: time.Minute,
		}

		err := p.Ensure(t.Context())
		require.NoError(t, err)

		p.lastRefreshedAt = time.Now().Add(-2 * time.Minute)

		err = p.Ensure(t.Context())
		require.NoError(t, err)

		assert.Equal(t, int32(2), calls.Load())
	})
}
