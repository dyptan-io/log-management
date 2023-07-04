# Log Management - A Snow Software coding test

- [Overview](#overview)
- [API Specification](#api-specification)
- [Development](#development)
- [Setup and Testing](#setup-and-testing)

## Overview

The implementation of Snow coding challenge for Log Management solution is consists of two main components: a web server
(Log receiver) and log watcher (Log shipper). Both components share a single [Dockerfile](Dockerfile) and [internal](internal)
libraries for convenience.

It utilizes some commonly used libraries and tools to facilitate development, such as:

- [oapi-codegen](https://github.com/deepmap/oapi-codegen) Go-centric OpenAPI Client and Server Code Generator.
- [go-chi](https://github.com/go-chi/chi) A lightweight, idiomatic Go HTTP server and router.
- [slog](https://pkg.go.dev/golang.org/x/exp/slog) A new structured logger from standard library.

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
docker-compose up --build
```

Add log files to configured directory (/testdata as default) and fetch collected logs:

```sh
curl http://localhost:8080/v1/logs
```

Example output:

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


```sh
curl http://localhost:8080/v1/logs/a5843dcb-9f21-4123-9c7c-688f0e8b88a7
```

Example output:

```json
{
  "attributes": {},
  "id": "a5843dcb-9f21-4123-9c7c-688f0e8b88a7",
  "message": "Task faulted: 'Failed to listen on port 80'",
  "severity": "Error",
  "timestamp": "2021-11-10T13:18:52Z"
}
```
