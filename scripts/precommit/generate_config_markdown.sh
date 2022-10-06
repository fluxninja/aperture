#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/docs >/dev/null
make generate-config-markdown
popd >/dev/null
