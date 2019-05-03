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
