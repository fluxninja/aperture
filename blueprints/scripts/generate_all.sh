#!/usr/bin/env bash
set -euo pipefail

script_root=$(dirname "$0")
blueprints_root=${script_root}/..

FIND="find"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

# for all subdirectories within "$blueprints_root"/lib containing config.libsonnet, generate README
$FIND "$blueprints_root"/lib -type f -name config.libsonnet | while read -r files; do
	dir=$(dirname "$files")
	echo "Generating README for $dir"
	python "${blueprints_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
done

# for all directories within "$blueprints_root"/bundles, generate examples
$FIND "$blueprints_root"/examples -mindepth 1 -maxdepth 1 -type d | while read -r dir; do
	echo "Generating examples for $dir"
	# save the contents of $dir/gen/policies/example.yaml for comparison
	old_example_yaml=$(cat "$dir"/gen/policies/example.yaml 2>/dev/null || true)

	# generate example blueprint
	python "${blueprints_root}"/scripts/generate-bundle.py --output "$dir"/gen/ \
		--config "$dir"/example.jsonnet

	npx prettier --write "$dir"/gen/.. || true

	new_example_yaml=$(cat "$dir"/gen/policies/example.yaml)

	if [[ "$old_example_yaml" != "$new_example_yaml" ]]; then
		mkdir -p "$dir"/gen/graph
		# fail if commands below fails
		go run -mod=mod "${blueprints_root}"/../cmd/circuit-compiler/main.go \
			-cr "$dir"/gen/policies/example.yaml \
			-dot "$dir"/gen/graph/graph.dot
		# if exit code is not 0 then remove example.yaml
		# shellcheck disable=SC2181
		if [[ $? -ne 0 ]]; then
			rm "$dir"/gen/policies/example.yaml
			# exit script with error
			exit 1
		fi
		dot -Tsvg "$dir"/gen/graph/graph.dot >"$dir"/gen/graph/graph.svg
	fi
done
