# Changelog

## v5.4.0 (Oct 1 2018)

* [logger] Use official lib of rollbar in rollbarplugin

## v5.3.1 (Sep 27 2018)

* [graceful] Fix ListenAndServeTLS to correctly use a tls.Listener

## v5.3.0 (Sep 26 2018)

* Add `graceful` package for graceful restart/shutdown of HTTP servers

## v5.2.1 (Sep 25 2018)

* Add `SkipLogSet` options to NSQ LB producer env init

## v5.2.0 (Jul 20 2018)

* Add `WhereQuery` in `mongo/document`

## v4.2.3 (Jun 14 2018)

* Add ability to choose nsq producing strategy from env NSQ_PRODUCER_STRATEGY 'fallback' or 'random', default is 'fallback'

## v4.2.2 (Jun 14 2018)

* New package `nsqlbproducer` to produce among multiple nsqd instances
  Several strategies:
  * Random (default)
  * Fallback: always produce on the first host of the list, and fallback on other if it fails

## v4.0.0 (Mar 13 2018)

* `mongo/document` Add SortFields on all Find/Where methods

  Example of usage:

```
document.Where(ctx, col, query, data, document.SortField("-created_at"))
document.Find(ctx, col, query, id, data, document.SortField("-created_at"))
```

## v3.0.0 / v3.1.0 (Mar 8 2018)

* Uniformize API for Find/Sort with Unscoped versions and Sort versions

## v2.1.0 (Mar 5 2018)

* Add `mongo/document` package for CRUD operations against the database and
  base structure to embed in models

  ```
  type Model struct {
    document.Base `bson:",inline"`
  }

  type ParanoidModel struct {
    document.Paranoid `bson:",inline"`
  }
  ```

## v2.0.1

* Update rollbar plugin configuration, only skip 4 levels of stack

## v2.0.0

* Change interface of nsqproducer.Producer, add Stop method

## v1.1.2

### global

* No more reference to go-internal-tools

## v1.1.1

### nsqconsumer

* Create nsqconsumer.Error to handle when no retry should be done

## v1.1.0

### errors

* New errors.Notef and errors.Wrapf with context to let the error handling system
  read the context and its logger

## v1.0.2

### logger

* New API of logrus removed `logger.AddHook(hook)`, it is now `logger.Hooks.Add(hook)`

## v1.0.1

### logger

* Fix plugin system, use pointer instead of value for manager

## v1.0.0

### mongo (api change)

* Add a logger in initialization
