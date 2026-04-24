package api

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	cryptoutils "github.com/Scalingo/go-utils/crypto"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

const (
	HeaderWebhookTimestamp  = "X-Link-Webhook-Timestamp"
	HeaderWebhookSignature  = "X-Link-Webhook-Signature"
	HeaderWebhookResourceID = "X-Link-Webhook-Resource-ID"
)

// ParseAndValidateWebhook decodes the payload and validates the webhook signature.
// The request body is restored before returning so callers can read it again if needed.
func ParseAndValidateWebhook(
	ctx context.Context,
	req *http.Request,
	secret string,
) (WebhookPluginStatusChangePayload, error) {
	if req == nil {
		return WebhookPluginStatusChangePayload{}, ErrWebhookRequestNil
	}

	log := logger.Get(ctx)

	if secret == "" {
		return WebhookPluginStatusChangePayload{}, ErrWebhookSecretMissing
	}

	if req.Body == nil {
		return WebhookPluginStatusChangePayload{}, ErrWebhookRequestBodyMissing
	}

	body, err := io.ReadAll(req.Body)
	closeErr := req.Body.Close()
	if closeErr != nil {
		log.WithError(closeErr).Error("Close webhook request body")
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return WebhookPluginStatusChangePayload{}, errors.Wrap(ctx, err, "read webhook body")
	}

	var payload WebhookPluginStatusChangePayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return WebhookPluginStatusChangePayload{}, errors.Wrap(ctx, ErrWebhookPayloadInvalid, err.Error())
	}

	timestamp := req.Header.Get(HeaderWebhookTimestamp)
	if timestamp == "" {
		return payload, ErrWebhookTimestampMissing
	}

	signature := req.Header.Get(HeaderWebhookSignature)
	if signature == "" {
		return payload, ErrWebhookSignatureMissing
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return payload, errors.Wrap(ctx, ErrWebhookSignatureInvalid, err.Error())
	}

	expectedSignature := cryptoutils.HMAC256([]byte(secret), []byte(timestamp+"."+string(body)))
	if subtle.ConstantTimeCompare(signatureBytes, expectedSignature) != 1 {
		return payload, ErrWebhookSignatureMismatch
	}

	return payload, nil
}
