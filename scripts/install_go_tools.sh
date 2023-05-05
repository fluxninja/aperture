#!/usr/bin/env bash

set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "${gitroot}" >/dev/null

./scripts/build_aperturectl.sh --force-build >/dev/null

if ! command -v go &>/dev/null; then
	printf 'Installing Go\n'
	./scripts/install_asdf_tools.sh setup golang
fi

printf 'Installing Go tools\n'

pushd "$gitroot"/tools/go >/dev/null

# first run go mod download
go mod download

tools=$(grep _ tools.go | awk -F'"' '{print $2}')

# use a parallel command to install go tools in parallel
parallel -j8 --no-notice --bar --eta go install ::: "$tools"

popd >/dev/null

if asdf where golang &>/dev/null; then
	asdf reshim golang
fi

popd >/dev/null
