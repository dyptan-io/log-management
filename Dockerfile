# syntax = docker/dockerfile:1.2
# https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/syntax.md

FROM golang:1.20-alpine3.17 AS builder

ARG package

WORKDIR /src/
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o srv -v -ldflags "-s -w" ./cmd/${package}

FROM alpine:3.17

COPY --from=builder /src/srv /usr/local/bin/

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
	CMD wget -q -O - http://localhost:8080/health || exit 1

ENTRYPOINT [ "srv" ]

