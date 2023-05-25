#!/usr/bin/env bash

set -e
git_root=$(git rev-parse --show-toplevel)

aperturectl="$git_root"/cmd/aperturectl/aperturectl
# check if aperturectl exists
if [ ! -f "$aperturectl" ]; then
	aperturectl="$("$git_root"/scripts/build_aperturectl.sh)"
fi

# function takes URI, name, values-file as input
function generate_compare() {
	local name=$1
	local values_file=$2
	local yaml1=$3
	local yaml2=$4

	"$aperturectl" blueprints generate \
		--uri "$git_root"/blueprints \
		--name "$name" \
		--values-file "$values_file" \
		--output-dir "tmp" \
		--no-yaml-modeline \
		--skip-pull \
		--overwrite

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
		echo "aperturectl generate did not match jsonnet library example (compared $yaml1" "and $yaml2))"
		echo "$comp"
		exit 1
	fi
	rm -rf tmp
}

# function generate_from_values takes name, values_file and output_dir as input and rendes the blueprint
function generate_from_values() {
	local name=$1
	local values_file=$2
	local output_dir=$3

	"$aperturectl" blueprints generate \
		--uri "$git_root"/blueprints \
		--name "$name" \
		--values-file "$values_file" \
		--output-dir "$output_dir" \
		--no-yaml-modeline \
		--skip-pull \
		--overwrite
}
