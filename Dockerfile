FROM golang:alpine

MAINTAINER Trevor Hartman <trevorhartman@gmail.com>

RUN apk add --update curl apache2-utils && rm -rf /var/cache/apk/*

COPY . /go/src/app

WORKDIR /go/src/app

ENTRYPOINT ["/usr/local/go/bin/go", "run", "main.go"]
