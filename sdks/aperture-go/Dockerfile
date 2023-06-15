# syntax=docker/dockerfile:1
FROM golang:1.20.4 AS builder

WORKDIR /src
COPY --link . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "aperture-go-example" "./example"

# Final image
FROM debian:bullseye-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  wget \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/aperture-go-example /local/bin/aperture-go-example

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
  CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD ["/local/bin/aperture-go-example"]
