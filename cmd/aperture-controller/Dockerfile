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
ARG APERTURECTL_BUILD_BUILDER_CONFIG_FILE
ENV APERTURECTL_BUILD_BUILDER_CONFIG_FILE=${APERTURECTL_BUILD_BUILDER_CONFIG_FILE}

# If config file exists, copy it to the build context
RUN if [[ -f "${APERTURECTL_BUILD_BUILDER_CONFIG_FILE}" ]]; then \
  cp "${APERTURECTL_BUILD_BUILDER_CONFIG_FILE}" ./config.yaml; \
  else \
  echo "APERTURECTL_BUILD_BUILDER_CONFIG_FILE is not set"; \
  touch ./config.yaml; \
  fi


RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=controller-1.19.6,sharing=private \
  --mount=type=cache,target=/root/.aperturectl,id=controller-1.19.6,sharing=private \
  /bin/bash -c \
  'go mod download; \
  go build -o cmd/aperturectl/aperturectl ./cmd/aperturectl; \
  ./cmd/aperturectl/aperturectl build controller -c ./config.yaml -o / --uri .; \
  '

# Final image
FROM debian:bullseye-slim AS controller

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --link --from=builder /aperture-controller /aperture-controller
CMD ["/aperture-controller"]
