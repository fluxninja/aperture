#!/usr/bin/env bash
set -euo pipefail

APERTURECTL_DIR=$(git rev-parse --show-toplevel)/cmd/aperturectl

cd "${APERTURECTL_DIR}"
if ! go build 1>&2; then
	exit 1
fi

echo -n "${APERTURECTL_DIR}/aperturectl"
