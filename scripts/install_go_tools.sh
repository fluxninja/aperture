#!/usr/bin/env bash

set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot"/tools/go >/dev/null

# first run go mod download
go mod download

tools=$(grep _ tools.go | awk -F'"' '{print $2}')

# use a parallel command to install go tools in parallel
parallel -j8 --no-notice --bar --eta go install ::: "$tools"

popd >/dev/null
