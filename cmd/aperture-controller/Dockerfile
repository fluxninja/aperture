# syntax=docker/dockerfile:1.4
FROM golang:1.20.4-bullseye AS builder

WORKDIR /src

COPY go.mod go.sum ./
# Cache go mod download as separate layer, but only in case of empty mount
# cache (like in CI).
RUN --mount=type=cache,target=/go/pkg/ \
  if [ ! -d "/go/pkg/mod" ]; then \
    export GOMODCACHE=/go/ci-pkg/mod && \
    go mod download && \
    cp -ar /go/ci-pkg/mod /go/pkg/mod; \
  fi

COPY --link . .

# App build
ARG APERTURECTL_BUILD_VERSION
ENV APERTURECTL_BUILD_VERSION=${APERTURECTL_BUILD_VERSION}
ARG APERTURECTL_BUILD_GIT_COMMIT_HASH
ENV APERTURECTL_BUILD_GIT_COMMIT_HASH=${APERTURECTL_BUILD_GIT_COMMIT_HASH}
ARG APERTURECTL_BUILD_GIT_BRANCH
ENV APERTURECTL_BUILD_GIT_BRANCH=${APERTURECTL_BUILD_GIT_BRANCH}
ARG APERTURECTL_BUILD_FLAGS
ENV APERTURECTL_BUILD_FLAGS=${APERTURECTL_BUILD_FLAGS}

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=controller-1.20.4,sharing=private \
  --mount=type=cache,target=/root/.aperturectl,id=controller-1.20.4,sharing=private \
  /bin/bash -c \
  'go mod download; \
  ./scripts/build_aperturectl.sh ./cmd/aperturectl; \
  ./cmd/aperturectl/aperturectl build controller -o / --uri .; \
  '

# Final image
FROM debian:bullseye-slim AS controller

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --link --from=builder /aperture-controller /aperture-controller

CMD ["/aperture-controller"]
