#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)

go run "$git_root"/cmd/aperturectl/main.go blueprints generate \
	--uri "$git_root"/blueprints \
	--name policies/static-rate-limiting \
	--values-file values.yaml \
	--output-dir "tmp" \
	--skip-pull

rm -rf tmp

# copy over raw values.yaml as well
cp "$git_root"/blueprints/policies/static-rate-limiting/values-required.yaml raw_values.yaml
