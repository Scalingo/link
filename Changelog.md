# To be released

# [2020-03-06] 1.8.3

* Make LeaseTime according to the IP KeepaliveInterval if set

# [2020-03-06] 1.8.2

* Update mocks for testing

# [2020-03-06] 1.8.1

* Update api.Client interface

# [2020-03-06] 1.8.0

* Feature: Allow the configuration of KeepaliveInterval and HealthcheckInterval by IP

# [2020-01-29] 1.7.1

* Dependencies: Update etcd and go-utils libraries

# [2019-12-17] 1.7.0

* Bugfix: when a save leased is no found, consider it expired to renew it completely
* ARP: Only send 3 Gratuitous ARP packets after becoming primary, configurable with env `ARP_GRATUITOUS_COUNT`

# [2019-11-08] 1.6.3

* Improve logging less useless errors when un-threatening healthcheck fail

# [2019-11-05] 1.6.2

* Improve logging after retries have been executed in etcd locker/healthcheck

# [2019-11-05] 1.6.1

* Add retry logics when refreshing ETCD lock (default 5 retries)
* Reduce logging verbosity when healthcheck is negative

# [2019-10-24] 1.6.0

* Add PPROF web
* Add BOOTING state to manage boot order.
* Fix: High CPU consumption when doing gratuitous ARP
* Fix: Do not remove IP on restart

# [2019-08-12] 1.5.2

* Fix: Logging is too verbose

# [2019-08-12] 1.5.1

* Do not loose IP on restart

# [2019-08-12] 1.5.0

* Release IP on standby
* Change LeaseDuration to 5 * KeepAliveInterval

# [2019-07-18] 1.4.3

* Update: logrus-rollbar to v1.3.1

# [2019-07-18] 1.4.2

* Fix: Health checker were stopped too early

# [2019-05-17] 1.4.1

* Fix regression: Remove IP on Stop

# [2019-05-15] 1.4.0

* Do not release IP on standby
* Do not use another lease if it has not expired

# [2019-05-14] 1.3.4

* Lease time is now 3 times the KeepaliveInterval

# [2019-05-03] 1.3.3

* Update go-philae to v4.3.2

# [2019-05-03] 1.3.2

* Update go-philae to v4.3.1
* Use go modules instead of dep

# [2019-04-30] 1.3.1

* Update go-philae to v4.3.0

# [2019-04-25] 1.3.0

* Do not fail on first healthcheck failure, add `FAIL_COUNT_BEFORE_FAILOVER`
  environment variable to configure the number of healthcheck failure before
  failover.

# [2019-04-15] 1.2.0

* Fix Client interface
* Make probes more verbose

# [2018-12-10] 1.1.0

* Release IP early if someone else got the lock
* Add the `version` endpoint and command

# [2018-11-29] 1.0.0

* First stable release
