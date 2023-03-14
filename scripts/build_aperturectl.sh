#!/usr/bin/env bash
set -eux

APERTURECTL_DIR=$(git rev-parse --show-toplevel)/cmd/aperturectl

APERTURECTL_BUILD_VERSION=${APERTURECTL_BUILD_VERSION:-0.0.1}
BUILD_TIME=${BUILD_TIME:-$(date -Iseconds)}
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
HOSTNAME=$(hostname)

APERTURECTL_BUILD_GIT_BRANCH=${APERTURECTL_BUILD_GIT_BRANCH:-$(git branch --show-current)}
APERTURECTL_BUILD_GIT_COMMIT_HASH=${APERTURECTL_BUILD_GIT_COMMIT_HASH:-$(git log -n1 --format=%H)}

LDFLAGS="\
    ${LDFLAGS:-} \
    -X 'github.com/fluxninja/aperture/pkg/info.Version=${APERTURECTL_BUILD_VERSION}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildHost=${HOSTNAME}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildOS=${GOOS}/${GOARCH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.BuildTime=${BUILD_TIME}' \
    -X 'github.com/fluxninja/aperture/pkg/info.GitBranch=${APERTURECTL_BUILD_GIT_BRANCH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.GitCommitHash=${APERTURECTL_BUILD_GIT_COMMIT_HASH}' \
    -X 'github.com/fluxninja/aperture/pkg/info.Service=aperturectl' \
    -X 'github.com/fluxninja/aperture/pkg/info.Prefix=aperture' \
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

echo -n "${APERTURECTL_DIR}/aperturectl"
