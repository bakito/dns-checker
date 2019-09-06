FROM golang:1.13 as builder

WORKDIR /go/src/github.com/bakito/dns-checker

COPY . /go/src/github.com/bakito/dns-checker/

RUN apt-get update && apt-get install -y xz-utils && \
  curl -SL --fail --silent --show-error https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz | tar --wildcards -xJ --strip-components 1 */upx

ENV GOPROXY=https://goproxy.io

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o dns-checker && \
   ./upx --ultra-brute -q dns-checker

# application image

FROM scratch

LABEL maintainer="bakito <github@bakito.ch>"
EXPOSE 2112
USER 1001
ENTRYPOINT ["/go/bin/dns-checker"]

COPY --from=builder /go/src/github.com/bakito/dns-checker/dns-checker /go/bin/dns-checker