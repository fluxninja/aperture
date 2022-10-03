# syntax=docker/dockerfile:1
FROM golang:1.19-buster AS builder

WORKDIR /src

COPY --link . .

# App build
RUN --mount=type=cache,target=/go/pkg/ \
    --mount=type=cache,target=/root/.cache/go-build/ <<-EOF
    go mod download
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /demoapp-go ./tools/demoapp-go
EOF

# Final image
FROM alpine:3.15.0
COPY --link --from=builder /demoapp-go /demoapp-go
ENTRYPOINT ["/demoapp-go"]
