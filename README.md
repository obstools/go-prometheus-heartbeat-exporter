# ![Prometheus Heartbeat Exporter - lightweight configurable multithreaded smokeping written on Golang](https://repository-images.githubusercontent.com/848341297/259630bf-4766-4265-86d2-6af5418f3299)

[![Go Report Card](https://goreportcard.com/badge/github.com/obstools/go-prometheus-heartbeat-exporter)](https://goreportcard.com/report/github.com/obstools/go-prometheus-heartbeat-exporter)
[![Codecov](https://codecov.io/gh/obstools/go-prometheus-heartbeat-exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/obstools/go-prometheus-heartbeat-exporter)
[![CircleCI](https://circleci.com/gh/obstools/go-prometheus-heartbeat-exporter/tree/master.svg?style=svg)](https://circleci.com/gh/obstools/go-prometheus-heartbeat-exporter/tree/master)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/obstools/go-prometheus-heartbeat-exporter)](https://github.com/obstools/go-prometheus-heartbeat-exporter/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/obstools/go-prometheus-heartbeat-exporter)](https://pkg.go.dev/github.com/obstools/go-prometheus-heartbeat-exporter)
[![GitHub](https://img.shields.io/github/license/obstools/go-prometheus-heartbeat-exporter)](LICENSE.txt)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

`heartbeat` is a lightweight configurable multithreaded Prometheus Heartbeat Exporter that helps you monitor your services with minimal overhead. Built for [Prometheus](https://prometheus.io) - the popular open-source monitoring and alerting system, it provides essential health metrics about your infrastructure.

Inspired by the principles of [smokeping](https://oss.oetiker.ch/smokeping/), `heartbeat` offers robust monitoring capabilities that allow you to measure the latency and availability of your services over time. With features like fast performance, easy setup, and efficient multi-target monitoring, `heartbeat` makes it simple to track service availability and response times.
Perfect for both small deployments and large-scale infrastructures, it seamlessly integrates into your existing Prometheus ecosystem while keeping resource usage low.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Supported connections](#supported-connections)
  - [Configuring](#configuring)
  - [Starting server](#starting-server)
    - [Available command line arguments](#available-command-line-arguments)
  - [Stopping server](#stopping-server)
- [Contributing](#contributing)
- [License](#license)
- [Code of Conduct](#code-of-conduct)
- [Credits](#credits)
- [Versioning](#versioning)
- [Changelog](CHANGELOG.md)

## Features

- Lightweight and efficient Prometheus exporter written in Go
- Multithreaded architecture for improved performance
- Configurable via YAML configuration file
- Graceful shutdown support
- Prometheus metrics exposure
- Comprehensive logging options
- Easy to deploy and configure

## Requirements

Golang 1.23+

## Installation

Install `heartbeat`:

```bash
go get github.com/obstools/go-prometheus-heartbeat-exporter
go install github.com/obstools/go-prometheus-heartbeat-exporter
```

Import `heartbeat` dependency into your code:

```go
package main

import heartbeat "github.com/obstools/go-prometheus-heartbeat-exporter"
```

## Usage

- [Supported connections](#supported-connections)
- [Configuring](#configuring)
- [Starting server](#starting-server)
  - [Available command line arguments](#available-command-line-arguments)
- [Stopping server](#stopping-server)

### Supported connections

Instances in the `heartbeat` configuration allow you to monitor various types of connections. Each instance can be configured with specific attributes such as `name`, `connection`, `url`, `query`, `interval`, and `timeout`. The `connection` attribute specifies the type of connection to be monitored, such as `postgres` for PostgreSQL, etc.

By providing a `query`, you can define specific operations to be executed on the database or other infrastructure elements, enabling you to check not only the connection status but also the performance of specific queries. This flexibility allows for comprehensive monitoring of your services and ensures that you can quickly identify and respond to issues as they arise.

> [!NOTE]
> Please make sure that your query is idempotent, as it can be executed multiple times during the `heartbeat` check. If `query` is not provided, `heartbeat` will check if connection is established only.

| Connection | Description | Query example |
| --- | --- | --- |
| `postgres` | Postgres database heartbeat. | `CREATE TABLE tmp (id SERIAL PRIMARY KEY); DROP TABLE tmp` |
| `redis` | Redis database heartbeat. Implemented only `SET`, `GET`, `DEL` commands. | `SET key1 value1; GET key1; DEL key1` |

### Configuring

`heartbeat` configuration is available as YAML file. You can also use environment variable interpolation in your configuration. This allows you to set sensitive information, such as database URLs or API keys, as environment variables and reference them in your YAML configuration. For example, you can set `url: '${DB_URL}'` and ensure that the `DB_URL` environment variable is defined in your environment. Available configuration options are described below.

> [!NOTE]
> Please keep in mind that instance name should be unique. Otherwise, `heartbeat` will fail to start.

| Configuration Key | Type     | Description | Example |
| --- | --- | --- | --- |
| `log_to_stdout` | optional | Enables logging to standard output. | `log_to_stdout: true` |
| `log_activity` | optional | Enables logging of heartbeat activity. | `log_activity: true` |
| `port` | required | Specifies Heartbeat Prometheus exporter server port number. | `port: 8080` |
| `metrics_route` | required | Defines the route for Prometheus exporter metrics. | `metrics_route: '/metrics'` |
| `instances` | required | List of instances to monitor. | `instances: [...]` |
| `name` | required | Name of the instance. Should be unique. | `name: 'postgres_1'` |
| `connection` | required | Type of connection to the instance. | `connection: 'postgres'` |
| `url` | required | Connection URL for the instance. | `url: 'postgres://localhost:5432/heartbeat_test'` |
| `query` | optional | Query to execute on the instance. | `query: 'CREATE TABLE tmp (id SERIAL PRIMARY KEY); DROP TABLE tmp'` |
| `interval` | required | Check interval for the instance in seconds. | `interval: 3` |
| `timeout` | required | Check timeout for the instance in seconds. | `timeout: 2` |

An example of the configuration file:

```yaml
# config.yml

---

log_to_stdout: true
log_activity: true
port: 8080
metrics_route: '/metrics'
instances:
  - name: 'postgres_1'
    connection: 'postgres'
    url: 'postgres://localhost:5432/heartbeat_test'
    query: 'CREATE TABLE tmp (id SERIAL PRIMARY KEY); DROP TABLE tmp'
    interval: 3
    timeout: 2
  - name: 'redis_1'
    connection: 'redis'
    url: 'redis://localhost:6379'
    query: 'SET key1 value1; GET key1; DEL key1'
    interval: 3
    timeout: 2
```

### Starting server

You can use `heartbeat` as binary. Just download the pre-compiled binary from the [releases page](https://github.com/obstools/go-prometheus-heartbeat-exporter/releases) and copy them to the desired location. For start server run command with path to the configuration file. You can use our bash script for automation this process like in the example below:

```bash
curl -sL https://raw.githubusercontent.com/obstools/go-prometheus-heartbeat-exporter/master/script/download.sh | bash
./heartbeat -config=config.yml
```

Passing environment variables to the `heartbeat:

```bash
SOME_ENV_VAR=123 ./heartbeat -config=config.yml
```

#### Available command line arguments

| Flag description | Example of usage |
| --- | --- |
| `-config` - path to the configuration file | `-config=config.yml` |
| `-v` - Prints current `heartbeat` binary build data. Doesn't run the server. | `-v` |

### Stopping server

`heartbeat` accepts 3 shutdown signals: `SIGINT`, `SIGQUIT`, `SIGTERM`.

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/obstools/go-prometheus-heartbeat-exporter>. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tickets](https://github.com/obstools/go-prometheus-heartbeat-exporter/issues). Be sure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

This golang package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the `heartbeat` project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/obstools/go-prometheus-heartbeat-exporter/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/obstools/go-prometheus-heartbeat-exporter/stargazers) for showing their support

## Versioning

`heartbeat` uses [Semantic Versioning 2.0.0](https://semver.org)
