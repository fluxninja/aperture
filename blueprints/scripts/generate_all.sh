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

# for all directories within "$blueprints_root"/blueprints, generate README
$FIND "$blueprints_root"/blueprints -mindepth 1 -maxdepth 1 -type d | while read -r dir; do
	python "${blueprints_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
	# save the contents of $dir/example/gen/policies/example.yaml for comparison
	old_example_yaml=$(cat "$dir"/example/gen/policies/example.yaml 2>/dev/null || true)

	# generate example blueprint
	python "${blueprints_root}"/scripts/aperture-generate.py --output "$dir"/example/gen/ \
		--config "$dir"/example/example.jsonnet

	npx prettier --write "$dir"/example/gen/.. || true

	new_example_yaml=$(cat "$dir"/example/gen/policies/example.yaml)

	if [[ "$old_example_yaml" != "$new_example_yaml" ]]; then
		mkdir -p "$dir"/example/gen/graph
		# fail if commands below fails
		go run -mod=mod "${blueprints_root}"/../cmd/circuit-compiler/main.go \
			-policy "$dir"/example/gen/policies/example.yaml \
			-dot "$dir"/example/gen/graph/graph.dot
		# if exit code is not 0 then remove example.yaml
		# shellcheck disable=SC2181
		if [[ $? -ne 0 ]]; then
			rm "$dir"/example/gen/policies/example.yaml
			# exit script with error
			exit 1
		fi
		dot -Tsvg "$dir"/example/gen/graph/graph.dot >"$dir"/example/gen/graph/graph.svg
	fi
done
