#!/usr/bin/env bash
set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null

make helm-lint

popd >/dev/null
