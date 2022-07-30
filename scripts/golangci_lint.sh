#!/usr/bin/env bash
set -euo pipefail

args=("${@}")
golangci-lint run --color=always "${args[@]}" ./...
code=$?

exit $code
