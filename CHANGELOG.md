# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2024-10-27

### Added

- Added ability to heartbeat redis instances

```yml
instances:
  - name: 'redis_1'
    connection: 'redis'
    url: 'redis://localhost:6379'
    query: 'SET key1 value1; GET key1; DEL key1'
    interval: 3
    timeout: 2
```

## [0.2.0] - 2024-10-26

### Added

- Added feature to specify custom instance query

## [0.1.2] - 2024-10-26

### Fixed

- Fixed issue with package namespace

## [0.1.1] - 2024-10-25

### Added

- Added download script

## [0.1.0] - 2024-10-25

### Added

- First release of `heartbeat`.
