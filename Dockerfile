FROM golang:alpine

MAINTAINER Trevor Hartman <trevorhartman@gmail.com>

RUN apk add --update curl apache2-utils && rm -rf /var/cache/apk/*

COPY . /go/

RUN go build -o /go/bin/main .

ENTRYPOINT ["main"]
