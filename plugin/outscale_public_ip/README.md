# Outscale Public IP Plugin

This plugin manages the Outscale Public IPs.
It assigns a Public IP to an Outscale Network Interface when the endpoint is activated.
On de-activation, it attempts to remove the Public IP from the Network Interface.

The Control Loop is run every minute by default and will re-assigns the Public IP to the NIC if it is not.

## Environment Variables

- `OUTSCALE_PUBLIC_IP_REFRESH_INTERVAL`: Interval between two control loop for an endpoint.

## JSON Configuration

| Name           | Type   | Optional | Description                                                                                                 |
| -------------- | ------ | -------- | ----------------------------------------------------------------------------------------------------------- |
| `access_key`   | string | no       | Outscale Access Key                                                                                         |
| `secret_key`   | string | no       | Outscale Secret Key                                                                                         |
| `region`       | string | no       | Outscale Region                                                                                             |
| `public_ip_id` | string | no       | ID of the Outscale Public IP to manage                                                                      |
| `nic_id`       | string | no       | ID of the Outscale Network Interface to which the Public IP will be assigned once the endpoint is activated |

### Example

```json
{
  "access_key": "YOUR_ACCESS_KEY",
  "secret_key": "YOUR_SECRET_KEY",
  "region": "eu-west-2",
  "public_ip_id": "eipalloc-12345678",
  "nic_id": "eni-12345678"
}
```

## Outscale EIM Configuration

The minimum EIM policy needed to use this plugin is:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["api:LinkPublicIp", "api:ReadPublicIps", "api:UnlinkPublicIp"],
      "Resource": ["*"]
    }
  ]
}
```
