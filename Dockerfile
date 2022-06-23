FROM golang:1.16-alpine3.15 as build
LABEL maintainer="https://github.com/nekochans"
WORKDIR /go/app
COPY go.* .
RUN go mod download
COPY . .
ARG AIR_VERSION=v1.40.2
RUN set -eux && \
  apk update && \
  apk add --no-cache git && \
  go install github.com/cosmtrek/air@${AIR_VERSION}
