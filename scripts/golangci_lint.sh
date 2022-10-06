#!/usr/bin/env bash
set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

args=("${@}")
golangci-lint run --color=always "${args[@]}" "$gitroot"/...
code=$?

exit $code
