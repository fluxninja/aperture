#!/usr/bin/env bash
set -euo pipefail

readarray -t module_dirs < <(go list -f '{{.Dir}}' -m)

parallel --eta --no-notice --line-buffer --tag --colsep ' ' 'cd {} && go test ./...' ::: "${module_dirs[@]}"
