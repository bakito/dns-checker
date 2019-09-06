#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -o dns-checker main.go
upx dns-checker