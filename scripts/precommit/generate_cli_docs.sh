#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/docs >/dev/null
rm -rf ./content/reference/aperture-cli/aperturectl*.md
go run ./tools/aperturectl/generate-docs.go
npx prettier --prose-wrap="preserve" ./content/reference/aperture-cli/ --write
popd >/dev/null
