FROM golang:1.12-alpine

RUN apk update; \
    apk add --no-cache git; \
    mkdir go-auth-manager;

ADD . /go/src/github.com/LucasFrezarini/go-auth-manager

WORKDIR /go/src/github.com/LucasFrezarini/go-auth-manager

ENV GO111MODULE on

RUN go mod download 
RUN go install server/server.go

ENTRYPOINT /go/bin/server

EXPOSE 8080