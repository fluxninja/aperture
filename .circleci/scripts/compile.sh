#!/usr/bin/env bash

set -euo pipefail
set -x

GIT_BRANCH="$(git branch --show-current)"
GIT_COMMIT_HASH="$(git log -n1 --format=%H)"
GOOS="$(go env GOOS)"
export GIT_BRANCH GIT_COMMIT_HASH GOOS
mkdir -p "$HOME/go"
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
export RACE=""

export CGO_ENABLED=1
export PREFIX=aperture
export LDFLAGS='-s -w -extldflags "-Wl,--allow-multiple-definition"'

: "${VERSION?VERSION needs to be set}"

case "${1:-}" in
  agent)
    SOURCE="./cmd/aperture-agent" TARGET="./dist/aperture-agent" ./pkg/info/build.sh
    for plugin_dir in ./plugins/*/aperture-plugin-*; do
      plugin="$(basename "${plugin_dir}")"
      echo "Building plugin ${plugin}"
      SOURCE="${plugin_dir}" TARGET="./dist/plugins/${plugin}.so" ./pkg/plugins/build.sh
    done
    ;;
  cli)
    SOURCE="./cmd/aperturectl" TARGET="./dist/aperturectl" ./pkg/info/build.sh
    ;;
  *)
    printf "UNKNOWN COMPONENT '%s' - valid are 'agent', 'cli'.\n" "${1:-}"
    exit 1
  ;;
esac
