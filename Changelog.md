# To be released

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
