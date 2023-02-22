#!/usr/bin/env bash
set -euo pipefail

echo Generating libsonnet library
git_root=$(git rev-parse --show-toplevel)
python "$git_root"/scripts/jsonnet-lib-gen.py --output-dir "$git_root"/blueprints/gen "$git_root"/docs/gen/policy/policy.yaml
tk fmt "$git_root"/blueprints/gen
npx prettier --write "$git_root"/blueprints/gen/jsonschema/*.json
git add "$git_root"/blueprints/gen

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

	python "${git_root}"/scripts/blueprint-assets-generator.py "$dir"

	gen_dir="$dir"/gen
	gen_files=("$gen_dir"/values.yaml "$gen_dir"/values-required.yaml "$gen_dir"/dynamic-config-values.yaml "$gen_dir"/dynamic-config-values-required.yaml "$gen_dir"/definitions.json "$gen_dir"/dynamic-config-definitions.json)
	for gen_file in "${gen_files[@]}"; do
		if [ -f "$gen_file" ]; then
			npx prettier --write "$gen_file"
		fi
	done

done

# run prettier on generated readme docs
$FIND "$git_root"/docs/content/reference/policies/bundled-blueprints -type f -name '*.md' | while read -r files; do
	npx prettier --write "$files"
done
