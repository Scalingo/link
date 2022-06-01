# Changelog

## To be released

* chore(deps): bump github.com/gofrs/uuid from 4.1.0+incompatible to 4.2.0+incompatible
* chore(deps): bump go.etcd.io/etcd/api/v3 from 3.5.0 to 3.5.4
* chore(deps): bump go.etcd.io/etcd/client/v3 from 3.5.0 to 3.5.4
* chore(deps): bump github.com/Scalingo/go-handlers from 1.4.0 to 1.4.3
* chore(deps): bump github.com/Scalingo/go-utils/errors from 1.0.0 to 1.1.0
* chore(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1
* chore(deps): bump github.com/urfave/cli from 1.22.5 to 1.22.9

## [2021-10-21] v2.0.2

* chore(deps): bump github.com/Scalingo/go-philae/v4 from 4.4.2 to 4.4.3
* chore(deps): bump github.com/Scalingo/go-utils/logger from 1.0.0 to v1.1.0

## [2021-08-23] v2.0.1

* chore(deps): bump github.com/olekukonko/tablewriter v0.0.0-20180912035003-be2c049b30cc => v0.0.5
* chore(deps): replace github.com/satori/go.uuid with github.com/gofrs/uuid
* chore(deps): update github.com/gofrs/uuid v3.4.0+incompatible => v4.0.0+incompatible
* chore(deps): update github.com/j-keck/arping v0.0.0-20160618110441-2cf9dc699c56 => v1.0.2
* chore(deps): update github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e => v3.0.0
* chore(deps): update github.com/looplab/fsm v0.0.0-20180515091235-f980bdb68a89 => v0.3.0
* chore(deps): bump go.etcd.io to Go Module version
* chore(deps): update github.com/Scalingo/go-utils/etcd v1.0.1 => v1.1.0
* chore: replace Travis CI with GitHub Actions to release new versions

## [2021-07-30] v2.0.0

* chore(Dependabot): Update various dependencies
* Bump github.com/sirupsen/logrus from 1.8.0 to 1.8.1
* Bump github.com/golang/mock from 1.5.0 to 1.6.0
* Add `/failover` route and `failover` command to force a failover on an ACTIVATED IP.
* Remove `/try-get-lock` route and `try-get-lock` command in favor of `failover`
* Destroying an IP is now synchronous
* Remove `KeepaliveInterval` from the IP model
* Add a way to update IP healthchecks
* Bad request body key is `error` instead of `msg`
* Not found body is `{"resource": "IP", "error" : "not found"}`

## [2020-11-20] v1.9.3

* Update deps: use github.com/Scalingo/go-philae/v4@v4.4.2

## [2020-11-19] v1.9.2

* Update deps: use go.etcd.io/etcd/v3 instead of github.com/coreos/etcd

## [2020-11-19] v1.9.1

* Update deps: github.com/Scalingo/go-utils, use submodules instead of global

## [2020-05-06] v1.9.0

* Use a single ETCD lease per server instead of one per IP to reduce load on the etcd cluster
* Validate the health check port value to be in the range [1:65535]

## [2020-03-06] v1.8.4

* Fix condition leading to extra lease renewal for short keepalive durations

## [2020-03-06] v1.8.3

* Make LeaseTime according to the IP KeepaliveInterval if set

## [2020-03-06] v1.8.2

* Update mocks for testing

## [2020-03-06] v1.8.1

* Update api.Client interface

## [2020-03-06] v1.8.0

* Feature: Allow the configuration of KeepaliveInterval and HealthcheckInterval by IP

## [2020-01-29] 1.7.1

* Dependencies: Update etcd and go-utils libraries

## [2019-12-17] 1.7.0

* Bugfix: when a save leased is no found, consider it expired to renew it completely
* ARP: Only send 3 Gratuitous ARP packets after becoming primary, configurable with env `ARP_GRATUITOUS_COUNT`

## [2019-11-08] 1.6.3

* Improve logging less useless errors when un-threatening healthcheck fail

## [2019-11-05] 1.6.2

* Improve logging after retries have been executed in etcd locker/healthcheck

## [2019-11-05] 1.6.1

* Add retry logics when refreshing ETCD lock (default 5 retries)
* Reduce logging verbosity when healthcheck is negative

## [2019-10-24] 1.6.0

* Add PPROF web
* Add BOOTING state to manage boot order.
* Fix: High CPU consumption when doing gratuitous ARP
* Fix: Do not remove IP on restart

## [2019-08-12] 1.5.2

* Fix: Logging is too verbose

## [2019-08-12] 1.5.1

* Do not loose IP on restart

## [2019-08-12] 1.5.0

* Release IP on standby
* Change LeaseDuration to 5 * KeepAliveInterval

## [2019-07-18] 1.4.3

* Update: logrus-rollbar to v1.3.1

## [2019-07-18] 1.4.2

* Fix: Health checker were stopped too early

## [2019-05-17] 1.4.1

* Fix regression: Remove IP on Stop

## [2019-05-15] 1.4.0

* Do not release IP on standby
* Do not use another lease if it has not expired

## [2019-05-14] 1.3.4

* Lease time is now 3 times the KeepaliveInterval

## [2019-05-03] 1.3.3

* Update go-philae to v4.3.2

## [2019-05-03] 1.3.2

* Update go-philae to v4.3.1
* Use go modules instead of dep

## [2019-04-30] 1.3.1

* Update go-philae to v4.3.0

## [2019-04-25] 1.3.0

* Do not fail on first healthcheck failure, add `FAIL_COUNT_BEFORE_FAILOVER`
  environment variable to configure the number of healthcheck failure before
  failover.

## [2019-04-15] 1.2.0

* Fix Client interface
* Make probes more verbose

## [2018-12-10] 1.1.0

* Release IP early if someone else got the lock
* Add the `version` endpoint and command

## [2018-11-29] 1.0.0

* First stable release
