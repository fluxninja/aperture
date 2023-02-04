#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/docs >/dev/null
go run ./tools/aperturectl/generate-docs.go
npx prettier --prose-wrap="preserve" ./content/reference/aperturectl/ --write
popd >/dev/null
