FROM golang:1.16-alpine3.15 as build
LABEL maintainer="https://github.com/nekochans"
WORKDIR /go/app
COPY go.* .
RUN go mod download
COPY . .
ARG AIR_VERSION=v1.40.2
ARG DLV_VERSION=v1.8.3
ENV CGO_ENABLED 0
RUN set -eux && \
  apk update && \
  apk add --no-cache git && \
  go install github.com/cosmtrek/air@${AIR_VERSION} && \
  go install github.com/go-delve/delve/cmd/dlv@${DLV_VERSION}
