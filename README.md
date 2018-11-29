# LinK
[![Build Status](https://travis-ci.org/Scalingo/link.svg?branch=master)](https://travis-ci.org/Scalingo/link)

> Link is not Keepalived

LinK is a networking agent that will let multiple hosts share a virtual IP. It
chooses which host must bind this IP and inform other members of the
network of the host owning this IP.

The IP owner election is performed using etcd lease system and other hosts on
this network is informed of the current IP owner using gratuitous ARP
requests (see [How do we bind IPs?](#how-do-we-bind-the-ips)).

To ease the cluster administration, LinK comes with it's
[own CLI](https://github.com/Scalingo/link/tree/master/cmd/link-client/).


## Demo

![demo](https://raw.githubusercontent.com/Scalingo/link/master/media/demo.gif)

## Project goals

1. KISS: our goal is to follow the UNIX philosophy: "Do one thing and do it
   well". This component is only responsible of the IP attribution part. It
   will not manage load balancing or other higher level stuff.
1. If an IP is registered on the cluster there must always be *at least one*
   server that binds the IP

## Architecture

** No central manager** Each agent only have knowledge of their local
configuration. They do not know nor care if other IP exists or if other hosts
have the same IP configured. The synchronization is done by creating locks in
etcd.

** Fault resilience** If for any reason something went wrong (lost connection
with etcd) LinK will always try to have **at least** one host this means that
if one agent fails to contact the etcd cluster it will take the IP.

## Installation

In order to be able to run LinK, you must have a working etcd cluster.
Installation and configuration instructions are available on the [etcd
website](https://coreos.com/etcd/docs/latest/getting-started-with-etcd.html).

> Link uses etcd v3 API. So you'll need etcd version 3.0.0 or higher.

The easiest way to get link up and running is to use pre-build binary available
on the [release pages](https://github.com/Scalingo/link/releases).

### Configuration

LinK configuration is entirely done by setting environment variables.

- `INTERFACE`: Name of the interface where LinK should add and remove IPs.
- `USER`: Username used for basic auth
- `PASSWORD`: Password used for basic auth
- `PORT` (default: 1313): Port where the LinK HTTP interface will be available
- `ETCD_HOSTS`: The different endpoints of etcd members
- `ETCD_TLS_CERT`: Path to the TLS X.509 certificate
- `ETCD_TLS_KEY`: Path to the private key authenticating the certificate
- `ETCD_CACERT`: Path to the CA cert signing the etcd member certificates

## Endpoints

- `GET /ips`: List all currently configured IPs
- `POST /ips`: Add an IP
- `GET /ips/:id`: Get a single IP
- `DELETE /ips/:id`: Remove an IP
- `POST /ips/:id/lock`: Try to get the lock on this IP

## How do we bind the IPs?

To add an interface LinK adds the IP to the configured interface and send an
unsolicited ARP request on the network (see [Gratuitous
ARP](https://wiki.wireshark.org/Gratuitous_ARP)).

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


## Dev environment

To make it work in dev you might want to some a dummy interfaces:

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
