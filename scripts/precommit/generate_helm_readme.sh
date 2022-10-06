#!/usr/bin/env bash
set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null

make generate-helm-readme

popd >/dev/null
