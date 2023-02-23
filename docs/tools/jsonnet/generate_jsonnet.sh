#!/bin/bash

gitroot=$(git rev-parse --show-toplevel)
export gitroot

scriptroot=$(dirname "$0")
export scriptroot

docsdir=$gitroot/docs
export docsdir

blueprints_root="${gitroot}/blueprints"
export blueprints_root

GREP="grep"
SED="sed"
AWK="awk"
FIND="find"
# check whether we are using macOS
if [ "$(uname)" == "Darwin" ]; then
	GREP="ggrep"
	SED="gsed"
	AWK="gawk"
	FIND="gfind"
fi
export GREP
export SED
export AWK
export FIND

# accept a --force flag to force the regeneration of all graphs
# check if $1 is bounded and equal to --force
if [[ -n "${1:-}" ]] && [[ "$1" == "--force" ]]; then
	force=true
else
	force=false
fi
export force

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

function generate_jsonnet() {
	set -euo pipefail
	dir=$(dirname "$1")
	filename=$(basename "$1")
	# make a unique tmp directory
	tmpdir=$(mktemp -d)

	filenameNoExt="${filename%.*}"
	out_dir="$dir"/assets/gen/"$filenameNoExt"/jsonnet
	rm -rf "$out_dir"/*.jsonnet || true

	#shellcheck disable=SC2002,SC2016
	cat "$1" | $SED -n '/```jsonnet/,/```/p' | $GREP -vP '^```$' >"$tmpdir"/records.txt
	# use awk to separate out jsonnet_records using RS='```' into an array of sections
	#shellcheck disable=SC2016
	$AWK -v tmpdir="$tmpdir" '{RS="```jsonnet"} NR > 1 { print $0 > tmpdir"/jsonnet_section_" ++i}' "$tmpdir"/records.txt

	# for each jsonnet section in tmp directory
	jsonnet_section_files=$($FIND "$tmpdir" -type f -name "jsonnet_section_*" | sort -n)
	count=0
	for jsonnet_section_file in $jsonnet_section_files; do
		echo "Processing $1 :: $count"
		# ignore if the jsonnet file contains "@include:"
		if $GREP -qP '@include:' "$jsonnet_section_file"; then
			continue
		fi

		# mkdir -p "$out_dir" if it doesn't exist
		if [ ! -d "$out_dir" ]; then
			mkdir -p "$out_dir"
		fi

		jsonnetfilepath="$out_dir"/"$filenameNoExt"_"$count".jsonnet
		mv "$jsonnet_section_file" "$jsonnetfilepath"
		# tanka fmt "$jsonnetfilepath"
		tk fmt "$jsonnetfilepath"

		# increment count
		count=$((count + 1))
	done
	rm -rf "$tmpdir"
}

export -f generate_jsonnet

parallel -j4 --no-notice --bar --eta generate_jsonnet ::: "$($FIND "$docsdir"/content -type f -name "*.md")"

function generate_jsonnet_files() {
	set -euo pipefail
	jsonnet_file="$1"
	echo "Processing $jsonnet_file"
	dir=$(dirname "$jsonnet_file")
	# remove extension and add .yaml
	yamlfilepath="${jsonnet_file%.*}".yaml
	jsonfilepath="${jsonnet_file%.*}".json
	# tmpdir
	tmpdir=$(mktemp -d)

	# cp jsonnet file to tmp file
	tmpjsonnetfilepath="$tmpdir"/"$(basename "$jsonnet_file")"
	cp "$jsonnet_file" "$tmpjsonnetfilepath"

	# replace github.com/fluxninja/aperture/blueprints with $"gitroot"/blueprints
	$SED -i "s|github.com/fluxninja/aperture/blueprints|$gitroot/blueprints|g" "$tmpjsonnetfilepath"
	jsonnet -J "$gitroot"/blueprints/vendor "$tmpjsonnetfilepath" >"$jsonfilepath"
	# if the file is a policy kind then generate mermaid diagram
	if [ "$(yq e '.kind == "Policy"' "$jsonfilepath")" = "true" ]; then
		old_yaml_file_contents=""
		if [ -f "$yamlfilepath" ]; then
			old_yaml_file_contents=$(cat "$yamlfilepath")
		fi
		# convert the policy to yaml
		go run "$scriptroot"/json2yaml.go "$jsonfilepath" "$yamlfilepath"
		rm -rf "$jsonfilepath"
		# run prettier
		npx prettier --write "$yamlfilepath"
		# generate mermaid diagram
		# compile the policy and generate mermaid if yaml has changed
		if [ "$old_yaml_file_contents" != "$(cat "$yamlfilepath")" ] || [ "$force" = true ]; then
			# generate mermaid diagram
			mermaidfilepath="${jsonnet_file%.*}".mmd
			# compile the policy
			go run "$gitroot"/cmd/aperturectl/main.go compile --cr "$yamlfilepath" --mermaid "$mermaidfilepath"
		else
			# still validate the policy with compiler
			go run "$gitroot"/cmd/aperturectl/main.go compile --cr "$yamlfilepath"
		fi
	else
		npx prettier --write "$jsonfilepath"
	fi
	rm -rf "$tmpdir"
}

export -f generate_jsonnet_files

# find all jsonnet files in docs/content directory and generate them
parallel -j4 --no-notice --bar --eta generate_jsonnet_files ::: "$($FIND "$docsdir"/content -type f -name "*.jsonnet")"
