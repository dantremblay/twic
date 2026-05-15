## Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY local/ local/
COPY vendor/ vendor/

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG GITSTATE=clean

RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w \
              -X github.com/kassisol/twic/version.Version=${VERSION} \
              -X github.com/kassisol/twic/version.GitCommit=${COMMIT} \
              -X github.com/kassisol/twic/version.GitState=${GITSTATE} \
              -X github.com/kassisol/twic/version.BuildDate=$(date +%s)" \
    -o /twic .

## Runtime stage
FROM alpine:3

ARG VERSION=dev

LABEL io.harbormaster.image.maintainer="hbm@kassisol.com"
LABEL io.harbormaster.image.version=$VERSION
LABEL io.harbormaster.image.description="TWIC is an application for managing certificates to connect to the Docker daemon using TLS"

RUN apk add --no-cache bash curl ca-certificates

COPY --from=builder /twic /usr/local/bin/twic
COPY scripts/image/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
