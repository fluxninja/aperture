# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS builder

WORKDIR /src
COPY --link . .

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/ \
  go mod download

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/ \
  CGO_ENABLED=0 go build \
    --ldflags "-s -w -extldflags \"-Wl,--allow-multiple-definition\"" \
    -a \
    -o aperture-operator \
  operator/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /src/aperture-operator .
USER 65532:65532

ENTRYPOINT ["/aperture-operator"]
