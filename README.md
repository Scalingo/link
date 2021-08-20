# LinK v2.0.0

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

> LinK uses etcd v3 API and makes use of `LeaseValue` comparison in transactions. Hence you need etcd version 3.3.0 or higher.

The easiest way to get LinK up and running is to use pre-build binary available
on the [release pages](https://github.com/Scalingo/link/releases).

## State machine

Each LinK agent can be in any of these three states:

- `ACTIVATED`: This machine owns the virtual IP
- `STANDBY`: This machine does not own the virtual IP but is available for election
- `FAILING`: Health checks for this host failed, this machine is not available for election
- `BOOTING`: The VIP just started to join the cluster and is waiting for an election

At any point five types of events can happen:
- `fault`: There was some error when coordinating with other nodes.
- `elected`: This machine was elected to own the virtual IP.
- `demoted`: This machine just lost ownership of the virtual IP.
- `health_check_fail`: The health checks configured with this IP failed.
- `health_check_success`: The health checks configured with this IP succeeded.


This is what the state machine looks like:

![LinK state machine](./state_machine.png)

## Configuration

LinK configuration is entirely done by setting environment variables.

- `INTERFACE`: Name of the interface where LinK should add and remove IPs.
- `HOSTNAME`: Name of the host.
- `USER`: Username used for basic auth
- `PASSWORD`: Password used for basic auth
- `PORT` (default: 1313): Port where the LinK HTTP interface will be available
- `KEEPALIVE_INTERVAL`: Duration of the lease given to a VIP. If a node is down, it can take up to KEEPALIVE_INTERVAL seconds to failover.
- `KEEPALIVE_RETRY`: Number of communication errors with etcd needed before considering the etcd cluster down.
- `HEALTH_CHECK_INTERVAL`: Interval between two health check queries.
- `HEALTH_CHECK_TIMEOUT`: Max duration of a health check.
- `FAIL_COUNT_BEFORE_FAILOVER`: Number of failed health checks needed before failing over.
- `ARP_GRATUITOUS_INTERVAL`: Time between two gratuitous ARP packets.
- `ARP_GRATUITOUS_COUNT`: Number of gratuitous ARP packets sent when an IP becomes ACTIVATED.
- `ETCD_HOSTS`: The different endpoints of etcd members
- `ETCD_TLS_CERT`: Path to the TLS X.509 certificate
- `ETCD_TLS_KEY`: Path to the private key authenticating the certificate
- `ETCD_CACERT`: Path to the CA cert signing the etcd member certificates

## Endpoints

- `GET /ips`: List all currently configured IPs
- `POST /ips`: Add an IP
- `GET /ips/:id`: Get a single IP
- `DELETE /ips/:id`: Remove an IP
- `POST /ips/:id/failover`: Trigger a failover on this IP (can only be launched on the master)

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

## Dev environment

To make it work in dev you might want to make some dummy interfaces:

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

## Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`

Commit, tag and create a new release:

```sh
git add CHANGELOG.md README.md
git commit -m "Bump v2.0.0"
git tag v2.0.0
git push origin master v2.0.0
hub release create v2.0.0
```
