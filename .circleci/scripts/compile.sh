#!/usr/bin/env bash

set -euo pipefail
set -x

APERTURECTL_BUILD_GIT_BRANCH="$(git branch --show-current)"
APERTURECTL_BUILD_GIT_COMMIT_HASH="$(git log -n1 --format=%H)"
GOOS="$(go env GOOS)"
export APERTURECTL_BUILD_GIT_BRANCH APERTURECTL_BUILD_GIT_COMMIT_HASH GOOS
mkdir -p "$HOME/go"
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"

: "${APERTURECTL_BUILD_VERSION?APERTURECTL_BUILD_VERSION needs to be set}"

case "${1:-}" in
agent)
	aperturectl="$(./scripts/build_aperturectl.sh)"
	"$aperturectl" build agent --output-dir ./dist --uri .
	;;
cli)
	aperturectl="$(./scripts/build_aperturectl.sh)"
	mkdir -p ./dist
	cp "$aperturectl" ./dist/aperturectl
	;;
*)
	printf "UNKNOWN COMPONENT '%s' - valid are 'agent', 'cli'.\n" "${1:-}"
	exit 1
	;;
esac
