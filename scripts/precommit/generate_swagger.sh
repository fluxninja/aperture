#!/usr/bin/env bash
set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null
# generate swagger
make go-generate-swagger

popd >/dev/null
