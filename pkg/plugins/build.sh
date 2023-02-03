#!/usr/bin/env bash
set -eux

# This script builds a Go plugin and injects build-time variables.

BUILD_TIME=$(date --rfc-3339=seconds)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)
PLUGIN_FILE=$(basename -- "${TARGET}")
PLUGIN="${PLUGIN_FILE%.*}"

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'main.Plugin=${PLUGIN}' \
    -X 'main.BuildHost=${HOSTNAME}' \
    -X 'main.BuildOS=${GOOS}/${GOARCH}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitBranch=${GIT_BRANCH}' \
    -X 'main.GitCommitHash=${GIT_COMMIT_HASH}' \
"

build_args=(
  -buildmode=plugin
  --ldflags "${LDFLAGS}"
  -o "${TARGET}"
  "${SOURCE}"
)
if [ -n "${RACE:-}" ]; then
  build_args+=( --race )
fi
go build "${build_args[@]}"
