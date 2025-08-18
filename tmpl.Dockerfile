FROM golang:1.25 as builder

WORKDIR /build

RUN apt-get update && apt-get install -y upx
COPY . .

ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH={{ .GoARCH }} \
    GOARM={{ .GoARM }}

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o dns-checker && \
    upx --ultra-brute -q dns-checker

# application image

FROM scratch

LABEL maintainer="bakito <github@bakito.ch>"
EXPOSE 2112
USER 1001
ENTRYPOINT ["/go/bin/dns-checker"]

COPY --from=builder /build/dns-checker /go/bin/dns-checker