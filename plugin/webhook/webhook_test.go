package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/models"
)

func TestPluginOnStatusChange(t *testing.T) {
	t.Run("activate sends payload and configured headers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "token-123", r.Header.Get("Authorization"))
			assert.Equal(t, "link", r.Header.Get("X-App"))

			var body map[string]any
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)

			assert.Equal(t, "vip-1", body["endpoint_id"])
			assert.Equal(t, Name, body["plugin"])
			assert.Equal(t, "ACTIVATED", body["status"])
			_, err = time.Parse(time.RFC3339Nano, body["changed_at"].(string))
			require.NoError(t, err)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint: models.Endpoint{ID: "vip-1", Plugin: Name},
			cfg: PluginConfig{
				URL: server.URL,
				Headers: map[string]string{
					"Authorization": "token-123",
					"X-App":         "link",
				},
			},
			httpClient: server.Client(),
		}

		err := p.Activate(context.Background())
		require.NoError(t, err)
	})

	t.Run("deactivate returns error on non 2xx status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		p := &Plugin{
			endpoint:   models.Endpoint{ID: "vip-2", Plugin: Name},
			cfg:        PluginConfig{URL: server.URL},
			httpClient: server.Client(),
		}

		err := p.Deactivate(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "webhook returned non-success status code")
	})
}
