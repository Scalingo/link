# LINK

> Link Is Not Keepalived

The goal of this project is to provide a simple and easy way to manage virtual
IPs. This project aims to be as KISS and dynamic as possible.

## Project promises

1. KISS our goal is not to rebuild Keepalived nor Pacemaker
1. If an IP is configured on a server there must always be *at least one* server that binds the IP

## How do we bind the IPs ?

To add an interface LINK will add the IP to the configured interface and sent a unsolicited ARP request on the network (see [Gratuitous ARP](https://wiki.wireshark.org/Gratuitous_ARP)).

This is the equivalent of:

```shell
ip addr add MY_IP dev MY_INTERFACE
arping -B -S MY_IP -I MY_INTERFACE
```

To unbind an IP we will just remove it from the interface.

This is the equivalent of:

```shell
ip addr del MY_IP dev MY_INTERFACE
```

## State machine

For each IP, the agent can either be `ACTIVATED` meaning that this machine owns the IP or `STANDBY` meaning that this machine does not own this IP yet.

At some point three types of events can happen:
- `fault`: There was some error when coordinating with other notes
- `elected`: This machine was elected to own the IP
- `demoted`: This machine just loosed ownership on the IP


This is what the state machine looks like:
![Sate Machine](./state_machine.png)

## Endpoints

- `GET /ips`: List all currently configured IPs
- `POST /ips`: Add an IP
- `DELETE /ips/:id`: Remove an IP
