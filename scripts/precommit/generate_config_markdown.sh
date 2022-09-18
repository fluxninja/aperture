#!/usr/bin/env bash
set -euo pipefail

pushd ./docs >/dev/null
make generate-config-markdown
popd >/dev/null
