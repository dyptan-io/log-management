# Log Management - A Snow Software coding test

- [Requirements](#requirements)
- [Overview](#overview)
- [API Specification](#api-specification)
- [Development](#development)
- [Setup and Testing](#setup-and-testing)

## Requirements

#### Log shipper

The log shipper should be able to read logs from a configurable source and send these forward to the receiver. The scope of this exercise is to
have the file system as a supported source and for this source, the following should to be supported:

- [x]  Detecting when new files are added to the observed directories on the file system
- [x]  Post the rows within the text files in the directories to the REST API exposed by the log receiver
- [x]  Delta updates is desirable but not required where added rows to existing log files are written after initial upload

#### Log receiver

The log receiver accepts incoming log rows, stores these and makes them available for an end-user to read.
An REST API should be exposed for receiving the log rows and for reading them

- [x] Only unique rows shall be stored
- [x] The id parameter in the incoming logs should be used for uniqueness
- [x] Basic filtering is a plus but not required where the user can filter on time received in the read API in the log receiver
- [x] The logs must be stored in-memory of the receiver
- [ ] A bonus (but not a requirement) is to store the logs in a persistent store in
addition to the in-memory store

## Overview

The implementation of Snow coding challenge for Log Management solution is consists of two main components: a web server
(Log receiver) and log watcher (Log shipper). Both components share a single [Dockerfile](Dockerfile) and [internal](internal)
libraries for convenience. Tests are added only for [InMemory](internal/platform/storage/inmemory_test.go) storage so far.

In addition to standard library, it utilizes some commonly used libraries and tools to facilitate development, such as:

- [oapi-codegen](github.com/oapi-codegen/oapi-codegen) Go-centric OpenAPI Client and Server Code Generator.
- [testify](https://github.com/stretchr/testify) A handy toolkit for assertion in tests.

## API Specification

For complete API documentation see the [OpenAPI Spec](api/v1.yaml) document.
You can view this document at [Swagger Editor](https://editor.swagger.io).

## Development

This project follows design and development principles described in:

- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)

## Setup and Testing

Use docker-compose to build and run both services:

```sh
docker compose up --build
```

Add log files to configured directory (/testdata as default) and fetch collected logs:

```sh
curl http://localhost:8080/v1/logs
```

```json
[
  {
    "attributes": {},
    "id": "a5843dcb-9f21-4123-9c7c-688f0e8b88a7",
    "message": "Task faulted: 'Failed to listen on port 80'",
    "severity": "Error",
    "timestamp": "2021-11-10T13:18:52Z"
  },
  {
    "attributes": {
      "test": "test"
    },
    "id": "0d3f329c-3d20-4975-9beb-cf4425d3a138",
    "message": "Task faulted 3: 'Failed to listen on port 80'",
    "severity": "Error",
    "timestamp": "2021-11-10T13:18:54Z"
  }
]
```

Get log entries by ID:

```sh
curl http://localhost:8080/v1/logs/a5843dcb-9f21-4123-9c7c-688f0e8b88a7
```

```json
{
  "attributes": {},
  "id": "a5843dcb-9f21-4123-9c7c-688f0e8b88a7",
  "message": "Task faulted: 'Failed to listen on port 80'",
  "severity": "Error",
  "timestamp": "2021-11-10T13:18:52Z"
}
```

Filter log entries by timestamp:

```sh
curl 'http://localhost:8080/v1/logs?from=2021-11-10T13%3A18%3A53Z&to=2021-11-10T13%3A18%3A55Z'
```

```json
[
  {
    "attributes": {
      "test": "test"
    },
    "id": "0d3f329c-3d20-4975-9beb-cf4425d3a138",
    "message": "Task faulted 3: 'Failed to listen on port 80'",
    "severity": "Error",
    "timestamp": "2021-11-10T13:18:54Z"
  }
]
```
