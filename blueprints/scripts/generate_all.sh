#!/usr/bin/env bash
set -euo pipefail

script_root=$(dirname "$0")
blueprints_root=${script_root}/..

FIND="find"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

# for all directories within "$blueprints_root"/blueprints, generate README
$FIND "$blueprints_root"/blueprints -mindepth 1 -maxdepth 1 -type d | while read -r dir; do
	python "${blueprints_root}"/scripts/blueprint-readme-generator.py "$dir"
	if which prettier >/dev/null 2>&1; then
		prettier --write "$dir"/README.md
	fi
	# generate example blueprint
	python "${blueprints_root}"/scripts/aperture-generate.py --output "$dir"/_gen/ \
		--config "$dir"/example.jsonnet
	go run -mod=mod "${blueprints_root}"/../cmd/circuit-compiler/main.go \
		-policy "$dir"/_gen/policies/example.yaml \
		-dot "$dir"/graph.dot
	dot -Tsvg "$dir"/graph.dot >"$dir"/graph.svg
	rm -rf "$dir"/_gen
done
