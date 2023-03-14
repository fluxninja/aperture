# syntax=docker/dockerfile:1.4
FROM golang:1.19.6-bullseye AS builder

WORKDIR /src
COPY --link . .

# App build
ARG APERTURE_VERSION
ENV APERTURE_VERSION=${APERTURE_VERSION}
ARG APERTURE_GIT_COMMIT_HASH
ENV APERTURE_GIT_COMMIT_HASH=${APERTURE_GIT_COMMIT_HASH}
ARG APERTURE_GIT_BRANCH
ENV APERTURE_GIT_BRANCH=${APERTURE_GIT_BRANCH}
ARG APERTURE_APERTURE_BUILDER_CONFIG_FILE
ENV APERTURE_BUILDER_CONFIG_FILE=${APERTURE_BUILDER_CONFIG_FILE}

# If config file exists, copy it to the build context
RUN if [[ -f "${APERTURE_BUILDER_CONFIG_FILE}" ]]; then \
  cp "${APERTURE_BUILDER_CONFIG_FILE}" ./config.yaml; \
  else \
  echo "APERTURE_BUILDER_CONFIG_FILE is not set"; \
  echo "bundled_extensions: [fluxninja, sentry]" > ./config.yaml; \
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
