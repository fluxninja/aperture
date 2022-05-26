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

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'aperture.tech/aperture/pkg/info.Service=${SERVICE}' \
    -X 'aperture.tech/aperture/pkg/info.BuildHost=${HOSTNAME}' \
    -X 'aperture.tech/aperture/pkg/info.BuildOS=${GOOS}/${GOARCH}' \
    -X 'aperture.tech/aperture/pkg/info.BuildTime=${BUILD_TIME}' \
    -X 'aperture.tech/aperture/pkg/info.GitBranch=${GIT_BRANCH}' \
    -X 'aperture.tech/aperture/pkg/info.GitCommitHash=${GIT_COMMIT_HASH}' \
    -X 'aperture.tech/aperture/pkg/info.Prefix=${PREFIX}' \
"
go build --ldflags "${LDFLAGS}" -o "${TARGET}" "${SOURCE}"
