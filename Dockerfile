FROM golang:1.16-alpine3.15 as build
LABEL maintainer="https://github.com/nekochans"
WORKDIR /go/app
COPY go.* .
RUN go mod download
COPY . .
RUN set -eux && \
  go build -o lgtm-cat-api ./cmd/local/main.go

FROM alpine:3.15
WORKDIR /app
COPY --from=build /go/app/lgtm-cat-api .
RUN set -x && \
  addgroup go && \
  adduser -D -G go go && \
  chown -R go:go /app/lgtm-cat-api
CMD ["./lgtm-cat-api"]
