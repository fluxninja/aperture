# syntax=docker/dockerfile:1.4
FROM golang:1.19.6-bullseye AS base

WORKDIR /src
COPY --link . .

RUN --mount=type=cache,target=/go/pkg/ \
  go mod download

FROM base AS builder

# App build
ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}
ARG GIT_BRANCH
ENV GIT_BRANCH=${GIT_BRANCH}
ARG VERSION
ENV VERSION=${VERSION}
ARG RACE
ENV RACE=${RACE}
ARG BUILD_TIME="$(date -Iseconds)"
ENV BUILD_TIME=${BUILD_TIME}

RUN  --mount=type=cache,target=/go/pkg/,readonly \
  --mount=type=cache,target=/root/.cache/go-build/,id=controller-1.19.6,sharing=private \
  /bin/bash -c \
  'set -euo pipefail; \
  shopt -s nullglob; \
  CGO_ENABLED=1 TARGET="/aperture-controller" PREFIX="aperture" SOURCE="./cmd/aperture-controller" \
  LDFLAGS="-extldflags \"-Wl,--allow-multiple-definition\"" \
  ./pkg/info/build.sh;\
  for plugin in ./plugins/{service,controller}/aperture-plugin-*; do\
  echo "building plugin $plugin";\
  CGO_ENABLED=1 TARGET="/plugins/$(basename $plugin).so" PREFIX="aperture" SOURCE="$plugin" \
  LDFLAGS="-s -w -extldflags \"-Wl,--allow-multiple-definition\"" \
  ./pkg/plugins/build.sh;\
  done\
  '

# Final image
FROM debian:bullseye-slim AS controller

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Ensure config dirs exists, even if not mounted
RUN mkdir -p /etc/aperture/aperture-controller/classifiers /etc/aperture/aperture-controller/policies /opt/aperture/aperture-controller/plugins

COPY --link --from=builder /aperture-controller /aperture-controller
CMD ["/aperture-controller"]

COPY --link --from=builder /plugins /opt/aperture/aperture-controller/plugins
