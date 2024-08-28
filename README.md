# ![Prometheus Heartbeat Exporter - lightweight configurable multithreaded smokeping written on Golang](https://repository-images.githubusercontent.com/42/42)

[![Go Report Card](https://goreportcard.com/badge/github.com/bestwebua/go-prometheus-heartbeat-exporter)](https://goreportcard.com/report/github.com/bestwebua/go-prometheus-heartbeat-exporter)
[![Codecov](https://codecov.io/gh/bestwebua/go-prometheus-heartbeat-exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/bestwebua/go-prometheus-heartbeat-exporter)
[![CircleCI](https://circleci.com/gh/bestwebua/go-prometheus-heartbeat-exporter/tree/master.svg?style=svg)](https://circleci.com/gh/bestwebua/go-prometheus-heartbeat-exporter/tree/master)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bestwebua/go-prometheus-heartbeat-exporter)](https://github.com/bestwebua/go-prometheus-heartbeat-exporter/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bestwebua/go-prometheus-heartbeat-exporter)](https://pkg.go.dev/github.com/bestwebua/go-prometheus-heartbeat-exporter)
[![GitHub](https://img.shields.io/github/license/bestwebua/go-prometheus-heartbeat-exporter)](LICENSE.txt)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

`heartbeat` is lightweight configurable multithreaded Prometheus Heartbeat Exporter.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Configuring with command line arguments](#configuring-with-command-line-arguments)
  - [Other options](#other-options)
  - [Stopping server](#stopping-server)
- [Contributing](#contributing)
- [License](#license)
- [Code of Conduct](#code-of-conduct)
- [Credits](#credits)
- [Versioning](#versioning)
- [Changelog](CHANGELOG.md)

## Features

- Configurable multithreaded Prometheus Heartbeat Exporter...

## Requirements

Golang 1.15+

## Installation

Install `heartbeat`:

```bash
go get github.com/bestwebua/go-prometheus-heartbeat-exporter
go install github.com/bestwebua/go-prometheus-heartbeat-exporter
```

Import `heartbeat` dependency into your code:

```go
package main

import heartbeat "github.com/bestwebua/go-prometheus-heartbeat-exporter"
```

## Usage

- [Configuring with command line arguments](#configuring-with-command-line-arguments)
- [Other options](#other-options)
- [Stopping server](#stopping-server)

### Configuring

You can use `heartbeat` as binary. Just download the pre-compiled binary from the [releases page](https://github.com/bestwebua/go-prometheus-heartbeat-exporter/releases) and copy them to the desired location. For start server run command with needed arguments. You can use our bash script for automation this process like in the example below:

```bash
curl -sL https://raw.githubusercontent.com/bestwebua/go-prometheus-heartbeat-exporter/master/script/download.sh | bash
./heartbeat -port=2525 -log
```

#### Configuring with command line arguments

`heartbeat` configuration is available as command line arguments specified in the list below:

| Flag description | Example of usage |
| --- | --- |
| `-port` - server port number. If not specified it will be assigned dynamically | `-port=8080` |

#### Other options

Available not configuration `heartbeat` options:

| Flag description | Example of usage |
| --- | --- |
| `-v` - Just prints current `heartbeat` binary build data (version, commit, datetime). Doesn't run the server. | `-v` |

#### Stopping server

`heartbeat` accepts 3 shutdown signals: `SIGINT`, `SIGQUIT`, `SIGTERM`.

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/bestwebua/go-prometheus-heartbeat-exporter>. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tickets](https://github.com/bestwebua/go-prometheus-heartbeat-exporter/issues). Be sure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

This golang package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the `heartbeat` project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/bestwebua/go-prometheus-heartbeat-exporter/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/bestwebua/go-prometheus-heartbeat-exporter/stargazers) for showing their support

## Versioning

`heartbeat` uses [Semantic Versioning 2.0.0](https://semver.org)
