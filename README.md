# LinK
[![Build Status](https://travis-ci.org/Scalingo/link.svg?branch=master)](https://travis-ci.org/Scalingo/link)

> Link is not Keepalived

The goal of this project is to provide a simple and easy way to manage virtual
IPs. This project aims to be as KISS and dynamic as possible.

## Project promises

1. KISS: our goal is not to rebuild Keepalived nor Pacemaker
1. If an IP is configured on a server there must always be *at least one* server that binds the IP

## How do we bind the IPs?

To add an interface LinK adds the IP to the configured interface and send an unsolicited ARP request on the network (see [Gratuitous ARP](https://wiki.wireshark.org/Gratuitous_ARP)).

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

Each IP can be in any of these three states:

- `ACTIVATED`: This machine owns the IP
- `STANDBY`: This machine does not own the IP but is available for election
- `FAILING`: Health checks for this IP failed, this machine is not available for election

At any point five types of events can happen:
- `fault`: There was some error when coordinating with other nodes
- `elected`: This machine was elected to own the IP
- `demoted`: This machine just loosed ownership on the IP
- `health_check_fail`: The health checks configured with this IP failed.
- `health_check_success`: The health checks configured with this IP succeeded.


This is what the state machine looks like:

![Sate Machine](./state_machine.png)

## Endpoints

- `GET /ips`: List all currently configured IPs
- `POST /ips`: Add an IP
- `GET /ips/:id`: Get a single IP
- `DELETE /ips/:id`: Remove an IP
- `POST /ips/:id/lock`: Force a link to try to get this IP

## Dev environment

To make it work in dev you might want to add a dummy interface:

```shell
modprobe dummy
ip link add eth10 type dummy
ip link set eth10 up
ip link add eth11 type dummy
ip link set eth11 up
ip link add eth12 type dummy
ip link set eth12 up
```

The script `start.sh` can be executed as root to automatically do that.
