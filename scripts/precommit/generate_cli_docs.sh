#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/docs >/dev/null
rm -rf ./content/get-started/aperture-cli/aperturectl*.md
go run ../cmd/aperturectl/gen-docs/generate-docs.go
npx prettier --prose-wrap="preserve" ./content/get-started/aperture-cli/ --write
popd >/dev/null
