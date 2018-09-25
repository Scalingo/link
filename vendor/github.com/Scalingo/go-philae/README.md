# Go Philae

Go Philae is the go implementation of our Philae health check protocol.

## Architecture

The go-philae library is based around the idea of probe. Each service has some dependencies. Those dependencies are checked using probes.

The prober is the component that will take all probes, check everyone of them, and return an aggregated result into a single response.

Finally, the handler is a small utility that will take a prober and an existing router and add an `/_health` endpoint to that router. This endpoint is configured to run check every time someone call it.

## Usage

To use it in an existing project, you will need to add a prober with some probes, pass it to the handler and generate the route.

```
router := handlers.NewRouter("http")
// ... configure your routes

probes := prober.NewProber()
probes.AddProbe(redisprobe.NewRedisProbeFromURL("redis", config.E["REDIS_URL"]))
/// ... add as many probes as needed

globalRouter := philaehandler.NewRouter(prober, router)

http.ListenAndServe(":8080", globalRouter)
```

## Creating a probe

To create a probe, you will need to implement the following interface:

```
type Probe interface {
  Name() string // Return the name of the probe
  Check() error // Return nil if the probe check was successfull or an error otherwise
}
```

That's all folks.

### Convention

* The check must be as lightweight as possible
* The check should not modify data used by the service
* The check should not depend on the service state
* The check should not take more than 3 second

## Existing probes

### DockerProbe

Check that the docker daemon is running.

#### Usage

```
dockerprobe.NewDockerProbe(name, endpoint string)
```

* name: The name of the probe
* endpoint: The docker API endpoint

### EtcdProbe

Check that a etcd server is running

#### Usage

```
etcdprobe.NewEtcdProbe(name string, client etcd.KeysAPI)
```

* name: The name of the probe
* client: An etcd Keys API client correctly configured

### GithubProbe

Check that GitHub isn't reporting any issue.
It use the official GitHub status API to check if there is no "major" problem with the GitHub infrastructure.

#### Usage

```
githubprobe.NewGithubProbe(name string)
```

* name: The name of the probe

### HTTPProbe

Check that an HTTP service is running fine.
It will send a GET request to an endpoint and check that the response code is in the 2XX or 3XX class.

#### Usage

```
httpprobe.NewHTTPProbe(name, endpoint sring, opts HTTPOptions)
```

* name: The name of the probe
* endpoint: Endpoint which should be checked (i.e.: http://google.com)
* opts: General Options

HTTPOptions params:

* Username used for basic auth
* Password used for basic auth
* Checker custom checker that will check the response sent by the server
* ExpectedStatusCode will check for a specific status code
* DialTimeout provide a custom timeout to first byte
* ResponseTimeout provide a custom timeout from first byte to the end of the response

### MongoProbe

Check that a MongoDB database is up and running.

#### Usage

```
mongoprobe.NewMongoProbe(name, url string)
```

* name: The name of the probe
* url: Url used to check the probe

### NsqProbe

Check that a nsq database is up and running.

#### Usage

```
nsqprobe.NewNSQProbe(name, host string, port int)
```

* name: The name of the probe
* host: the IP address (or FQDN) of the nsq server
* port: The port on which the nsq server is running

### PhilaeProbe

Check that another service using Philae probe is running and healthy.

#### Usage

```
philaeprobe.NewPhilaeProbe(name, endpoint string, dialTimeout, responseTimeout int)
```

* name: The name of the probe
* endpoint: The philae endpoint (i.e.: "http://test.dev/_health")
* dialTimeout, responseTimeout: see HTTPProbe

### RedisProbe

Check that a Redis server is up and running

#### Usage

```
redisprobe.NewRedisProbe(name, host, password) string
```

* name: The name of the probe
* host: The Redis host
* password: The password needed to access the database


```
redisprobe.NewRedisProbeFromURL(name, url string)
```

* name: The name of the probe
* url: The url of the Redis server (i.e.: "redis://:password@server.org")

### SampleProbe

A probe only used for testing. It will always return the same result

#### Usage

```
sampleprobe.NewSampleProbe(name string, result bool)
```

* name: The name of the probe
* result: Is the check successful or not

```
sampleprobe.NewTimedSampleProbe(name string, result bool, time time.Duration)
```

* name: The name of the probe
* result: Is the check successful or not
* time: The time the probe will take before returning a result

### StatusIOProbe

This probe will check that a service using StatusIO is healthy

#### Usage

```
statusioprobe.NewStatusIOProbe(name, id string)
```

* name: The name of the probe
* id: The StatusIO service id

### SwiftProbe

Check that a Swift host is up and healthy.
It creates a new connection and try to authenticate.

#### Usage

```
swiftprobe.NewSwiftProbe(name, url, region, tenant, username, password string)
```

* name: The name of the probe
* url: Url of the Swift host
* tenant: The tenant name needed to authenticate
* username: The username needed to authenticate
* password: The password needed to authenticate

### TCPProbe

Check that a TCP server accept connection.

#### Usage

```golang
tcpprobe.NewTCPProbe(name, endpoint string, opts TCPOptions)
```

* name: The name of the probe
* endpoint: Endpoint to probe

Options:

* Timeout (default 5s)

