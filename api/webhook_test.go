package api

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cryptoutils "github.com/Scalingo/go-utils/crypto"
)

func TestValidateWebhookRequest(t *testing.T) {
	payload := WebhookPluginStatusChangePayload{
		EndpointID: "vip-1",
		ResourceID: "resource-123",
		Plugin:     PluginWebhook,
		Status:     Activated,
	}

	body := []byte(`{"endpoint_id":"vip-1","resource_id":"resource-123","plugin":"webhook","status":"ACTIVATED"}`)

	t.Run("validate matching signature", func(t *testing.T) {
		req := newWebhookRequest(t, body)

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, payload, gotPayload)

		bodyAfterValidation, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.Equal(t, body, bodyAfterValidation)
	})

	t.Run("return payload on signature mismatch", func(t *testing.T) {
		req := newWebhookRequest(t, body)

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "wrong-secret")

		require.NoError(t, err)
		assert.False(t, valid)
		assert.Equal(t, payload, gotPayload)
	})

	t.Run("reject nil request", func(t *testing.T) {
		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), nil, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookRequestNil)
		assert.False(t, valid)
		assert.Equal(t, WebhookPluginStatusChangePayload{}, gotPayload)
	})

	t.Run("reject missing secret", func(t *testing.T) {
		req := newWebhookRequest(t, body)

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "")

		require.ErrorIs(t, err, ErrWebhookSecretMissing)
		assert.False(t, valid)
		assert.Equal(t, WebhookPluginStatusChangePayload{}, gotPayload)
	})

	t.Run("reject missing body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://example.com/webhook", nil)
		require.NoError(t, err)
		req.Body = nil

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookRequestBodyMissing)
		assert.False(t, valid)
		assert.Equal(t, WebhookPluginStatusChangePayload{}, gotPayload)
	})

	t.Run("reject invalid payload", func(t *testing.T) {
		req := newWebhookRequest(t, []byte(`{"endpoint_id"`))

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookPayloadInvalid)
		assert.False(t, valid)
		assert.Equal(t, WebhookPluginStatusChangePayload{}, gotPayload)
	})

	t.Run("reject missing timestamp", func(t *testing.T) {
		req := newWebhookRequest(t, body)
		req.Header.Del(HeaderWebhookTimestamp)

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookTimestampMissing)
		assert.False(t, valid)
		assert.Equal(t, payload, gotPayload)
	})

	t.Run("reject missing signature", func(t *testing.T) {
		req := newWebhookRequest(t, body)
		req.Header.Del(HeaderWebhookSignature)

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookSignatureMissing)
		assert.False(t, valid)
		assert.Equal(t, payload, gotPayload)
	})

	t.Run("reject invalid signature encoding", func(t *testing.T) {
		req := newWebhookRequest(t, body)
		req.Header.Set(HeaderWebhookSignature, "not-hex")

		valid, gotPayload, err := ValidateWebhookRequest(t.Context(), req, "shared-secret")

		require.ErrorIs(t, err, ErrWebhookSignatureInvalid)
		assert.False(t, valid)
		assert.Equal(t, payload, gotPayload)
	})
}

const webhookTimestamp = "1713867072"

func newWebhookRequest(t *testing.T, body []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "https://example.com/webhook", bytes.NewReader(body))
	require.NoError(t, err)

	req.Header.Set(HeaderWebhookTimestamp, webhookTimestamp)
	req.Header.Set(
		HeaderWebhookSignature,
		hex.EncodeToString(cryptoutils.HMAC256([]byte("shared-secret"), []byte(webhookTimestamp+"."+string(body)))),
	)

	return req
}
