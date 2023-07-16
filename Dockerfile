FROM golang:1.20.4 AS builder

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /cache
COPY go.mod /cache

RUN go mod download

WORKDIR /release
COPY . /release

RUN make linux

FROM dbscale/base-centos7.9.2009:2-amd64

ENV LANG C.UTF-8

COPY --from=builder /release/dist/config-wrapper /usr/local/bin/config-wrapper
COPY --from=builder /release/etc/config.toml /etc/config-wrapper/config.toml

CMD ["config-wrapper","-f","/etc/config-wrapper/config.toml"]