#!/usr/bin/env bash
set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot"/docs >/dev/null
make generate-jsonnet
make generate-mermaid
popd >/dev/null
