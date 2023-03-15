# syntax=docker/dockerfile:1.4
FROM golang:1.19.6-bullseye AS builder

WORKDIR /src
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

ENV CONFIG_FILE=cmd/aperture-agent/build-config.yaml

# if build-config.yaml doesn't exist, create empty build-config.yaml
RUN if [ ! -f "${CONFIG_FILE}" ]; then touch "${CONFIG_FILE}"; fi

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=agent-1.19.6,sharing=private \
  --mount=type=cache,target=/root/.aperturectl,id=agent-1.19.6,sharing=private \
  /bin/bash -c \
  'go mod download; \
  go build -o cmd/aperturectl/aperturectl ./cmd/aperturectl; \
  ./cmd/aperturectl/aperturectl build agent -o / --uri . -c "${CONFIG_FILE}"; \
  '

# Final image
FROM debian:bullseye-slim AS agent

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --link --from=builder /aperture-agent /aperture-agent

CMD ["/aperture-agent"]
