#!/usr/bin/env bash
set -eux

# This script builds a Go plugin and injects build-time variables.

VERSION=${VERSION:-0.0.1}
BUILD_TIME=${BUILD_TIME:-$(date -Iseconds)}
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)
PREFIX="${PREFIX:-aperture}"

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

if [ -n "${RACE:-}" ]; then
	build_args=(-race)
fi

build_args+=(
	-buildmode=plugin
	--ldflags "${LDFLAGS}"
	-o "${TARGET}"
	"${SOURCE}"
)

go build "${build_args[@]}"
