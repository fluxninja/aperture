# syntax=docker/dockerfile:1.4
FROM golang:1.20.6-bullseye AS builder

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

RUN --mount=type=cache,target=/go/pkg/ \
  --mount=type=cache,target=/root/.cache/go-build/,id=agent-1.20.6,sharing=private \
  --mount=type=cache,target=/root/.aperturectl,id=agent-1.20.6,sharing=private \
  /bin/bash -c \
  'set -eu; \
  go mod download; \
  ./scripts/build_aperturectl.sh ./cmd/aperturectl; \
  ./scripts/generate_full_build_config.sh ./extensions > build.yaml; \
  ./cmd/aperturectl/aperturectl build agent -c build.yaml -o / --uri .; \
  '

# Final image
FROM debian:bullseye-slim AS agent

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --link --from=builder /aperture-agent /aperture-agent

ARG USERNAME=aperture-agent
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && chown -R $USERNAME:$USERNAME /aperture-agent

USER $USERNAME

CMD ["/aperture-agent"]
