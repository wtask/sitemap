# Distribute smgen as docker image
FROM golang:1.12 as builder

LABEL maintainer="<Mike Mikhaylov, webtask@gmail.com>"

ENV \
	# golang env
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /src
COPY ./ ./

RUN	update-ca-certificates --fresh \
	&& go mod tidy \
	&& go mod download \
	&& mkdir /build \
	&& go build -a -ldflags '-extldflags "-static"' -o /build/smgen ./cmd/smgen/.

FROM scratch

USER smgen

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# copy application binary
COPY --from=builder /build/smgen ./app/

WORKDIR /data
VOLUME /data

ENTRYPOINT [ "/app/smgen" ]
CMD ["-output-dir=/data"]
