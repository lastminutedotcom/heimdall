FROM golang:1.12-stretch as bin

LABEL author="Team SRE <sre@lastminute.com>"

ARG APPNAME
COPY . /work
WORKDIR /work

RUN go build -a -ldflags '-extldflags "-static"' -o /$APPNAME .

FROM debian:stretch

COPY --from=bin /$APPNAME /heimdall

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y ca-certificates && update-ca-certificates

ENTRYPOINT ["/heimdall"]

