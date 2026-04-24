package api

import "errors"

var (
	ErrWebhookRequestNil         = errors.New("webhook request is required")
	ErrWebhookSecretMissing      = errors.New("webhook secret is required")
	ErrWebhookRequestBodyMissing = errors.New("webhook request body is required")
	ErrWebhookTimestampMissing   = errors.New("missing X-Link-Webhook-Timestamp header")
	ErrWebhookSignatureMissing   = errors.New("missing X-Link-Webhook-Signature header")
	ErrWebhookSignatureInvalid   = errors.New("invalid X-Link-Webhook-Signature header")
	ErrWebhookSignatureMismatch  = errors.New("webhook signature does not match")
	ErrWebhookPayloadInvalid     = errors.New("invalid webhook payload")
)

type ErrNotFound struct {
	message string
}

func (e ErrNotFound) Error() string {
	return e.message
}
