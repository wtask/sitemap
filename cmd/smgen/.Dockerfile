# Distribute smgen as docker image
FROM golang:1.12 as builder

LABEL maintainer="<Mike Mikhaylov, webtask@gmail.com>"

ENV \
	# golang env
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /src
VOLUME /src

RUN go mod download \
	&& mkdir /app \
	&& go build -a -ldflags '-extldflags "-static"' -o /app/smgen ./cmd/smgen/.
