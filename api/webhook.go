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
	HeaderWebhookTimestamp = "X-Link-Webhook-Timestamp"
	HeaderWebhookSignature = "X-Link-Webhook-Signature"
)

// ValidateWebhookRequest validates the webhook signature and decodes the payload.
// It returns whether the signature matches, the decoded payload, and an error for malformed requests.
// The request body is restored before returning so callers can read it again if needed.
func ValidateWebhookRequest(
	ctx context.Context,
	req *http.Request,
	secret string,
) (bool, WebhookPluginStatusChangePayload, error) {
	if req == nil {
		return false, WebhookPluginStatusChangePayload{}, ErrWebhookRequestNil
	}

	log := logger.Get(ctx)

	if secret == "" {
		return false, WebhookPluginStatusChangePayload{}, ErrWebhookSecretMissing
	}

	if req.Body == nil {
		return false, WebhookPluginStatusChangePayload{}, ErrWebhookRequestBodyMissing
	}

	body, err := io.ReadAll(req.Body)
	closeErr := req.Body.Close()
	if closeErr != nil {
		log.WithError(closeErr).Error("fail to close webhook request body")
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return false, WebhookPluginStatusChangePayload{}, errors.Wrap(ctx, err, "read webhook body")
	}

	var payload WebhookPluginStatusChangePayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return false, WebhookPluginStatusChangePayload{}, errors.Wrap(ctx, ErrWebhookPayloadInvalid, err.Error())
	}

	timestamp := req.Header.Get(HeaderWebhookTimestamp)
	if timestamp == "" {
		return false, payload, ErrWebhookTimestampMissing
	}

	signature := req.Header.Get(HeaderWebhookSignature)
	if signature == "" {
		return false, payload, ErrWebhookSignatureMissing
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, payload, errors.Wrap(ctx, ErrWebhookSignatureInvalid, err.Error())
	}

	expectedSignature := cryptoutils.HMAC256([]byte(secret), []byte(timestamp+"."+string(body)))

	return subtle.ConstantTimeCompare(signatureBytes, expectedSignature) == 1, payload, nil
}
