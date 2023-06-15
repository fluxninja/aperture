#!/usr/bin/env bash
set -euo pipefail

APERTURECTL_DIR=${1:-$(git rev-parse --show-toplevel)/cmd/aperturectl}
APERTURECTL_BINARY="${APERTURECTL_DIR}/aperturectl"

APERTURECTL_BUILD_VERSION=${APERTURECTL_BUILD_VERSION:-0.0.1}
BUILD_TIME=${BUILD_TIME:-$(date -Iseconds)}
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)

if git rev-parse --git-dir >/dev/null 2>&1; then
	APERTURECTL_BUILD_GIT_BRANCH=${APERTURECTL_BUILD_GIT_BRANCH:-$(git branch --show-current)}
	APERTURECTL_BUILD_GIT_COMMIT_HASH=${APERTURECTL_BUILD_GIT_COMMIT_HASH:-$(git log -n1 --format=%H)}
fi

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.Version=${APERTURECTL_BUILD_VERSION}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.BuildHost=${HOSTNAME}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.BuildOS=${GOOS}/${GOARCH}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.BuildTime=${BUILD_TIME}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.GitBranch=${APERTURECTL_BUILD_GIT_BRANCH}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.GitCommitHash=${APERTURECTL_BUILD_GIT_COMMIT_HASH}' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.Service=aperturectl' \
    -X 'github.com/fluxninja/aperture/v2/pkg/info.Prefix=aperture' \
"

if [ -n "${RACE:-}" ]; then
	build_args=(-race)
fi

build_args+=(
	--ldflags "${LDFLAGS}"
)

pushd "${APERTURECTL_DIR}" >/dev/null

if ! go build "${build_args[@]}" 1>&2; then
	exit 1
fi

popd >/dev/null

echo -n "${APERTURECTL_BINARY}"
