#!/usr/bin/env bash
set -euo pipefail

echo Generating libsonnet library
git_root=$(git rev-parse --show-toplevel)
export git_root

python "$git_root"/scripts/jsonnet-lib-gen.py --output-dir "$git_root"/blueprints/gen "$git_root"/docs/gen/policy/policy.yaml
tk fmt "$git_root"/blueprints/gen
npx prettier --write "$git_root"/blueprints/gen/jsonschema/*.json
git add "$git_root"/blueprints/gen

blueprints_root=${git_root}/blueprints
export blueprints_root

FIND="find"
if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi
export FIND

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

function generate_readme() {
	set -euo pipefail
	dir=$(dirname "$1")
	echo "Generating README and Sample Values for $dir"

	python "${blueprints_root}"/blueprint-assets-generator.py "$dir"

	gen_dir="$dir"/gen
	gen_files=("$gen_dir"/values.yaml "$gen_dir"/values-required.yaml "$gen_dir"/dynamic-config-values.yaml "$gen_dir"/dynamic-config-values-required.yaml "$gen_dir"/definitions.json "$gen_dir"/dynamic-config-definitions.json)
	for gen_file in "${gen_files[@]}"; do
		if [ -f "$gen_file" ]; then
			npx prettier --write "$gen_file"
		fi
	done
}

export -f generate_readme

parallel -j4 --no-notice --bar --eta generate_readme ::: "$($FIND "$blueprints_root" -type f -name config.libsonnet)"

# run prettier on generated readme docs
parallel -j4 --no-notice --bar --eta npx prettier --write ::: "$($FIND "$git_root"/docs/content/reference/policies/bundled-blueprints -type f -name '*.md')"
