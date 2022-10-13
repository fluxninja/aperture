# syntax=docker/dockerfile:1

FROM golang:1.19-buster AS builder

WORKDIR /src
COPY --link . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "aperture-go-example" "./example"

# Final image
FROM alpine:3.15.0

COPY --from=builder /src/aperture-go-example /local/bin/aperture-go-example

HEALTHCHECK CMD curl --fail http://localhost:8080/super || exit 1

ENTRYPOINT ["/local/bin/aperture-go-example"]
