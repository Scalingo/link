# ARP Plugin

This plugin manages IPs and announces them on the local network using ARP.

## Environment Variables

- `INTERFACE`: Name of the interface where LinK should add and remove IPs.
- `ARP_GRATUITOUS_COUNT`: Number of gratuitous ARP packets sent when an IP becomes ACTIVATED.

## JSON Configuration

| Name | Type   | Optional | Description                              |
| ---- | ------ | -------- | ---------------------------------------- |
| `ip` | string | no       | IP address to manage using CIDR notation |

### Example

```json
{
  "ip": "10.20.30.40/32"
}
```

## How do we bind the IPs?

To add an interface, LinK adds the IP to the configured interface and send an
unsolicited ARP request on the network (see [Gratuitous
ARP](https://wiki.wireshark.org/Gratuitous_ARP)).

This is the equivalent of:

```shell
ip addr add MY_IP dev MY_INTERFACE
arping -B -S MY_IP -I MY_INTERFACE
```

To unbind an IP, LinK removes it from the interface.

This is the equivalent of:

```shell
ip addr del MY_IP dev MY_INTERFACE
```
