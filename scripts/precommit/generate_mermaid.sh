#!/usr/bin/env bash
set -euo pipefail

pushd ./docs >/dev/null
make generate-mermaid
popd >/dev/null
