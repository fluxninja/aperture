#!/usr/bin/env bash
set -eux

# This script builds a Go plugin and injects build-time variables.

# Create BUILD_TIME if it doesn't exist
BUILD_TIME=${BUILD_TIME:-$(date -Iseconds)}
VERSION=${VERSION:-0.0.1}
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)
PLUGIN_FILE=$(basename -- "${TARGET}")
PLUGIN="${PLUGIN_FILE%.*}"
PREFIX="${PREFIX:-aperture}"

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'main.Plugin=${PLUGIN}' \
    -X 'main.BuildHost=${HOSTNAME}' \
    -X 'main.BuildOS=${GOOS}/${GOARCH}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitBranch=${GIT_BRANCH}' \
    -X 'main.GitCommitHash=${GIT_COMMIT_HASH}' \
    -X 'main.Version=${VERSION}' \
    -X 'main.Prefix=${PREFIX}' \
"

if [ -n "${RACE:-}" ]; then
  build_args=( -race )
fi

build_args+=(
  -buildmode=plugin
  --ldflags "${LDFLAGS}"
  -o "${TARGET}"
  "${SOURCE}"
)

go build "${build_args[@]}"
