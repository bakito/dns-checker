# Multi-stage build with explicit platform specification
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /build

ARG VERSION=main
ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /go/src/app

RUN apk add --no-cache upx ca-certificates tzdata

ENV CGO_ENABLED=0 \
  GOOS=linux

# Copy only module files first to maximize layer caching for deps.
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Now copy the rest of the source code.
COPY . .

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X github.com/bakito/dns-checker/version.Version=${VERSION}" \
    -o dns-checker .

RUN upx -q dns-checker

# application image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL maintainer="bakito <github@bakito.ch>"

RUN microdnf install bind-utils nc && \
    microdnf clean all
EXPOSE 2112
USER 1001
ENTRYPOINT ["/go/bin/dns-checker"]

COPY --from=builder /go/src/app/dns-checker /go/bin/dns-checker
