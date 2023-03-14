# syntax=docker/dockerfile:1.4
FROM golang:1.19.6-bullseye AS builder

WORKDIR /src
COPY --link . .

RUN --mount=type=cache,target=/go/pkg/ \
  go mod download

# build aperturectl
RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=agent-1.19.6,sharing=private \
  go build -o cmd/aperturectl/aperturectl ./cmd/aperturectl

# App build
ARG APERTURE_GIT_COMMIT_HASH
ENV APERTURE_GIT_COMMIT_HASH=${APERTURE_GIT_COMMIT_HASH}
ARG APERTURE_GIT_BRANCH
ENV APERTURE_GIT_BRANCH=${APERTURE_GIT_BRANCH}
ARG CONFIG_FILE
ENV CONFIG_FILE=${CONFIG_FILE}

# If config file exists, copy it to the build context
RUN if [[ -f "${CONFIG_FILE}" ]]; then \
  cp "${CONFIG_FILE}" ./config.yaml; \
  else \
  echo "CONFIG_FILE is not set"; \
  echo "bundled_extensions: [fluxninja, sentry]" > ./config.yaml; \
  fi


RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=agent-1.19.6,sharing=private \
  --mount=type=cache,target=/root/.aperturectl,sharing=private \
  ./cmd/aperturectl/aperturectl build agent -c ./config.yaml -o . --uri .

ENTRYPOINT ["tail", "-f", "/dev/null"]

# # Final image
# FROM debian:bullseye-slim AS agent

# RUN apt-get update \
#   && apt-get install -y --no-install-recommends ca-certificates \
#   && apt-get clean \
#   && rm -rf /var/lib/apt/lists/*

# COPY --link --from=builder /aperture-agent /aperture-agent
# CMD ["/aperture-agent"]
