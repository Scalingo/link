# Webhook Plugin

The webhook plugin sends an HTTP request when the endpoint status changes.

## Plugin config

```json
{
  "url": "https://example.com/link-events",
  "headers": {
    "Authorization": "Bearer my-token",
    "X-Custom": "custom-value"
  }
}
```

- `url` is required and must use `http` or `https`.
- `headers` is optional and contains extra headers injected in the request.

## Request payload

The plugin sends a `POST` request to `url` with a JSON body:

```json
{
  "endpoint_id": "vip-...",
  "plugin": "webhook",
  "previous_status": "STANDBY",
  "status": "ACTIVATED",
  "changed_at": "2026-03-27T12:34:56Z"
}
```
