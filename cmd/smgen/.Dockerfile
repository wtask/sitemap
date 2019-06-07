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

RUN go mod download \
	&& mkdir /build \
	&& go build -a -ldflags '-extldflags "-static"' -o /build/smgen ./cmd/smgen/.

FROM scratch
WORKDIR /app
COPY --from=builder /build/smgen ./
WORKDIR /data
VOLUME [ "/data" ]

ENTRYPOINT [ "/app/smgen" ]
CMD ["-output-dir=/data", "-num-workers=4"]