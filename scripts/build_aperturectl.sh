#!/usr/bin/env bash
set -euo pipefail

FORCE_BUILD=false
for arg in "$@"; do
	case $arg in
	--force-build)
		FORCE_BUILD=true
		shift
		;;
	esac
done

APERTURECTL_DIR=${1:-$(git rev-parse --show-toplevel)/cmd/aperturectl}
APERTURECTL_BINARY="${APERTURECTL_DIR}/aperturectl"

if [ -f "${APERTURECTL_BINARY}" ] && [ "$FORCE_BUILD" = false ]; then
	echo -n "${APERTURECTL_BINARY}"
	exit 0
fi

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

echo -n "${APERTURECTL_BINARY}"
