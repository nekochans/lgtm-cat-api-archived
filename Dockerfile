# syntax = docker/dockerfile:1.3

FROM golang:1.16-alpine3.15 as base
LABEL maintainer="https://github.com/nekochans"
WORKDIR /go/app
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download

FROM base AS dev
ARG AIR_VERSION=v1.40.2
ARG DLV_VERSION=v1.8.3
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  set -eux && \
  apk update && \
  apk add --no-cache git && \
  go install github.com/cosmtrek/air@${AIR_VERSION} && \
  go install github.com/go-delve/delve/cmd/dlv@${DLV_VERSION}

FROM base AS unit-test
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build
