#!/usr/bin/env bash
set -euo pipefail

git_root=$(git rev-parse --show-toplevel)
blueprints_root=${git_root}/blueprints

FIND="find"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

# for all subdirectories within "$blueprints_root"/lib containing config.libsonnet, generate README
$FIND "$blueprints_root" -type f -name config.libsonnet | while read -r files; do
	dir=$(dirname "$files")
	echo "Generating README and Sample Values for $dir"
	python "${git_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
	# if values.yaml exists, format it
	if [ -f "$dir"/values.yaml ]; then
		npx prettier --write "$dir"/values.yaml
	fi
	# if values_required.yaml exists, format it
	if [ -f "$dir"/values_required.yaml ]; then
		npx prettier --write "$dir"/values_required.yaml
	fi
done

$FIND "$git_root"/docs/content/reference/policies/bundled-blueprints -type f -name '*.md' | while read -r files; do
	npx prettier --write "$files"
done
