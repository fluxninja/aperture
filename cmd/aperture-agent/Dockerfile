# syntax=docker/dockerfile:1
FROM golang:1.19-buster AS base

WORKDIR /src
COPY --link . .

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/ \
  go mod download

# App build
ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}
ARG GIT_BRANCH
ENV GIT_BRANCH=${GIT_BRANCH}
ARG VERSION
ENV VERSION=${VERSION}

FROM base AS builder
RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/ \
  CGO_ENABLED=1 TARGET="/aperture-agent" PREFIX="aperture" SOURCE="./cmd/aperture-agent" LDFLAGS="-s -w" \
  ./pkg/info/build.sh


# BUILD PLUGINS STAGE
FROM base AS plugins-builder
# Plugins build
RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/ \
  /bin/bash -c 'shopt -s nullglob; for plugin in ./plugins/{service,agent}/aperture-plugin-*; do\
  echo "building plugin $plugin";\
  CGO_ENABLED=1 TARGET="/plugins/$(basename $plugin).so" PREFIX="aperture" SOURCE="$plugin" LDFLAGS="-s -w" \
  ./pkg/plugins/build.sh;\
  done\
  '


# Final image
FROM alpine:3.15.0
# Needed for linker to be compatible
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache curl
# Ensure config dirs exists, even if not mounted
RUN mkdir -p /etc/aperture/aperture-agent/flowcontrol /etc/aperture/aperture-agent/classifiers /etc/aperture/aperture-agent/plugins
COPY --link --from=builder /aperture-agent /aperture-agent
COPY --link --from=plugins-builder /plugins /var/lib/aperture/aperture-agent/plugins
CMD ["/aperture-agent"]
