#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)
aperturectl="$git_root"/cmd/aperturectl/aperturectl

"$aperturectl" blueprints generate \
	--uri "$git_root"/blueprints \
	--name policies/rate-limiting \
	--values-file values.yaml \
	--output-dir "tmp" \
	--skip-pull \
	--overwrite

rm -rf tmp

# copy over raw values.yaml as well
cp "$git_root"/blueprints/policies/rate-limiting/gen/values.yaml raw_values.yaml
