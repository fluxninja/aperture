#!/usr/bin/env sh
set -eux

# This script builds a Go binary and injects build-time variables.
# TODO: it should be used in every Dockerfile requiring version endpoint.

BUILD_TIME=$(date --rfc-3339=seconds)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)
SERVICE_FILE=$(basename -- "${TARGET}")
SERVICE="${SERVICE_FILE%.*}"

if [ -d .git ]; then
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    GIT_COMMIT_HASH=$(git rev-parse HEAD)
fi

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'github.com/fluxninja/aperture/pkg/info.Service=${SERVICE}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildHost=${HOSTNAME}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildOS=${GOOS}/${GOARCH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildTime=${BUILD_TIME}' \
    -X 'github.com/fluxninja/aperture/pkg/info.GitBranch=${GIT_BRANCH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.GitCommitHash=${GIT_COMMIT_HASH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.Prefix=${PREFIX}' \
"
go build --ldflags "${LDFLAGS}" -o "${TARGET}" "${SOURCE}"
