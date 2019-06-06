# Distribute smgen as docker image
FROM golang:1.12

LABEL maintainer="<Mike Mikhaylov, webtask@gmail.com>"

ENV \
	# golang env
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64
