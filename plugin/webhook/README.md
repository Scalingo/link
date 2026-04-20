# Webhook Plugin

The webhook plugin sends an HTTP request when the endpoint status changes.

## Plugin config

```json
{
  "url": "https://example.com/link-events",
  "resource_id": "resource-123",
  "secret": "shared-secret",
  "headers": {
    "Authorization": "Bearer my-token",
    "X-Custom": "custom-value"
  }
}
```

- `url` is required and must use `http` or `https`.
- `resource_id` is required and identifies the external resource tied to the webhook.
- `secret` is required, stored encrypted, and used to sign webhook requests.
- `headers` is optional and contains extra headers injected in the request.

## Request payload

The plugin sends a `POST` request to `url` with a JSON body:

```json
{
  "endpoint_id": "vip-...",
  "resource_id": "resource-123",
  "plugin": "webhook",
  "status": "ACTIVATED",
  "changed_at": "2026-03-27T12:34:56Z"
}
```

Requests include these authentication headers:

- `X-Link-Webhook-Timestamp`: Unix timestamp in seconds.
- `X-Link-Webhook-Signature`: hex-encoded HMAC-SHA256 of `"<timestamp>.<body>"`.
