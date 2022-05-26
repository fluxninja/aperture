#!/usr/bin/env bash
set -euo pipefail

args=("${@}")
readarray -t module_dirs < <(go list -f '{{.Dir}}' -m)
parallel --bar --eta --no-notice --line-buffer --tag --colsep ' ' "golangci-lint run --color=always ${args[*]} {}/..." ::: "${module_dirs[@]}"
code=$?

exit $code
