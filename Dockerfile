FROM golang:1.24-bullseye as builder

ARG VERSION=main
WORKDIR /build

RUN apt-get update && apt-get install -y upx
COPY . .

ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -a -installsuffix cgo -ldflags="-w -s -X github.com/bakito/dns-checker/version.Version=${VERSION}" -o dns-checker && \
    upx -q dns-checker

# application image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL maintainer="bakito <github@bakito.ch>"

RUN microdnf install bind-utils nc && \
    microdnf clean all
EXPOSE 2112
USER 1001
ENTRYPOINT ["/go/bin/dns-checker"]

COPY --from=builder /build/dns-checker /go/bin/dns-checker
