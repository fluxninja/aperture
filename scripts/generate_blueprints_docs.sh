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

	values_files=("$dir"/values.yaml "$dir"/values-required.yaml "$dir"/dynamic-config-values.yaml "$dir"/dynamic-config-values-required.yaml)
	for values_file in "${values_files[@]}"; do
		if [ -f "$values_file" ]; then
			npx prettier --write "$values_file"
		fi
	done

done

$FIND "$git_root"/docs/content/reference/policies/bundled-blueprints -type f -name '*.md' | while read -r files; do
	npx prettier --write "$files"
done
