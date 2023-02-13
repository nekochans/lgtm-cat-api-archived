# syntax = docker/dockerfile:1.3

FROM golang:1.19-bullseye as base
LABEL maintainer="https://github.com/nekochans"
WORKDIR /go/app
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download

FROM base AS dev
ARG AIR_VERSION=v1.40.2
ARG DLV_VERSION=v1.20.1
ARG MOQ_VERSION=v0.3.0
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  set -eux && \
  apt-get update && \
  apt-get install git && \
  go install github.com/cosmtrek/air@${AIR_VERSION} && \
  go install github.com/go-delve/delve/cmd/dlv@${DLV_VERSION} && \
  go install github.com/matryer/moq@${MOQ_VERSION}

FROM base AS unit-test
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build

FROM base AS build
ARG COMMIT_HASH
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-X main.release=${COMMIT_HASH} -s -w" -o /out/lgtm-cat-api .

FROM debian:bullseye-slim as production
COPY --from=build /out/lgtm-cat-api /
RUN set -x && \
  apt-get update &&  \
  apt-get install -y ca-certificates && \
  useradd go && \
  chown -R go:go /lgtm-cat-api
USER go
CMD ["./lgtm-cat-api"]
