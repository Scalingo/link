# Changelog

## v4.0.0

* Migrate to context embedded logger

## v3.1.0

* tcpprobe: add tcpprobe

## v3.0.1

* nsqprobe: use https scheme if tls config is set

## v3.0.0

* nsqprobe: add tls.Config argument to handle https test

## v2.0.1

* Sirupsen -> sirupsen

## v2.0.0

* Add ability to configure logger and verbosity of philaehandler

Breaking Change:

```
philaehandler.NewPhilaeRouter(http.Router, *prober.Prober)
philaehandler.NewHandler(*prober.Prober)

// become

philaehandler.NewPhilaeRouter(http.Router, *prober.Prober, philaehandler.HandlerOpts)
philaehandler.NewHandler(*prober.Prober, philaehandler.HandlerOpts)
```

## v1.0.0

* Initial release (semver version|to use with dep)
