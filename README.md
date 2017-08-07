# PromCache

An caching proxy for the Prometheus query endpoint for speeding up Grafana
dashboards.

## Description

This repository contains a small HTTP proxy server to run alongside your
Prometheus servers. By default it will listen on port 9091 and proxy all requests
to `localhost:9090`. Any request staring with `/api/v1/query` will receive
special threatment and get possibly cached. All other requests are passed
unaltered.

All requests to `/api/v1/query` will be forcefully cached. The `Cache-Control`
headers are removed, the `start` and `end` parameters will be rounded to the
next full minute and then passed to the upstream server.

This proxy is tailored for one specific use-case: Multiple users and/or dashboards
accessing the same Prometheus timeseries through Grafana. This can considerably
reduce the load on your Prometheus server while still providing _good enough_
data for most dashboards.

Just add this proxy as an additional data source in your _templated_ Grafana
dashboards and you can easily switch between direct and cached access to compare
speed and quality.

## Quickstart

Clone, build, start and point Grafana to `localhost:9091` instead of `localhost:9090`
(if Prometheus is running on localhost).

```
go build
./promcache
```

## Flags

| **Flag** | **Type** | Description |
| -------- | -------- | ----------- |
| `listen` | string   | Listen address. Default: `:9091`
| `dir`    | string   | Cache dir. Implies `-disk`
| `disk`   | bool     | Use disk cache. Default: `false`
| `verbose` | bool    | Verbose output
| `dumphttp` | bool   | Dump HTTP Request and Response headers
| `upstream` | string | Upstream host to connect to
| `duration` | int    | Forced cache duration. Default: `60`

## License

MIT
