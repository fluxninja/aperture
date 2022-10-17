# syntax=docker/dockerfile:1

FROM golang:1.19-alpine3.16 AS builder

WORKDIR /src
COPY --link . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "aperture-go-example" "./example"

# Final image
FROM alpine:3.16

COPY --from=builder /src/aperture-go-example /local/bin/aperture-go-example

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
    CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD ["/local/bin/aperture-go-example"]
