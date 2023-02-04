#!/usr/bin/env bash

set -e
git_root=$(git rev-parse --show-toplevel)

# function takes URI, name, values-file as input
function generate_compare() {
	local uri=$1
	local name=$2
	local values_file=$3
	local yaml1=$4
	local yaml2=$5

	go run "$git_root"/cmd/aperturectl/main.go blueprints generate \
		--uri "$uri" \
		--name "$name" \
		--values-file "$values_file" \
		--output-dir "tmp"

	# make sure yaml1 and yaml2 exist
	if [ ! -f "$yaml1" ]; then
		echo "yaml1 does not exist: $yaml1"
		exit 1
	fi
	if [ ! -f "$yaml2" ]; then
		echo "yaml2 does not exist: $yaml2"
		exit 1
	fi

	set +e
	# compare the generated yaml with the expected yaml using yq
	comp=$(diff <(yq -P 'sort_keys(..)' "$yaml1") <(yq -P 'sort_keys(..)' "$yaml2"))
	set -e
	if [ -n "$comp" ]; then
		echo "aperturectl generate did not match jsonnet library example"
		echo "$comp"
		exit 1
	fi
	rm -rf tmp
}
