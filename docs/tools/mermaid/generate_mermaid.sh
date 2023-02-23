#!/bin/bash

gitroot=$(git rev-parse --show-toplevel)
export gitroot

docsdir=$gitroot/docs
export docsdir

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

# accept a --force flag to force generation of all mermaid files
if [[ -n "${1:-}" ]] && [[ "$1" == "--force" ]]; then
	force=true
else
	force=false
fi
export force

# script that parses all markdown (md) files recursively in a directory
# and extracts content for the mermaid sections that begin with ```mermaid and end with ```
# and generates mmd files for each mermaid section
function generate_mermaid() {
	set -euo pipefail
	f="$1"
	dir=$(dirname "$f")
	filename=$(basename "$f")
	filenameNoExt="${filename%.*}"
	out_dir="$dir"/assets/gen/"$filenameNoExt"
	rm -f "$out_dir"/*.mmd
	tmpdir=$(mktemp -d)
	# extract all mermaid multiline sections that start with ```mermaid and end with ``` into an array of sections
	#shellcheck disable=SC2002,SC2016
	cat "$f" | $SED -n '/```mermaid/,/```/p' | $GREP -vP '^```$' >"$tmpdir"/records.txt
	# use awk to separate out mermaid_records using RS='```' into an array of sections
	#shellcheck disable=SC2016
	$AWK -v tmpdir="$tmpdir" '{RS="```mermaid"} NR > 1 { print $0 > tmpdir"/mermaid_section_" ++i}' "$tmpdir"/records.txt
	# find mermaid_section_* and sort them
	mermaid_section_files=$($FIND "$tmpdir" -type f -name "mermaid_section_*" | sort -n)
	count=0
	for mermaid_section_file in $mermaid_section_files; do
		# skip this file if it contains "@include:"
		if $GREP -q "@include:" "$mermaid_section_file"; then
			continue
		fi

		# mkdir -p "$out_dir" if it doesn't exist
		if [ ! -d "$out_dir" ]; then
			mkdir -p "$out_dir"
		fi

		# search for name in the comment - "%% name: <name>"
		# if found, use the name as the mmd file name
		name=$($GREP -P '^%% name: ' "$mermaid_section_file" | $SED -e 's/%% name: //')
		if [ -n "$name" ]; then
			outfilename="$name.mmd"
		else
			outfilename="$filename"_$count.mmd
		fi
		# generate mmd
		echo "generating $outfilename"
		mv "$mermaid_section_file" "$out_dir"/"$outfilename"
		# increment count
		count=$((count + 1))
	done
	rm -rf "$tmpdir"
}

export -f generate_mermaid

parallel -j4 --no-notice --bar --eta generate_mermaid ::: "$($FIND "$docsdir"/content -type f -name "*.md")"

# find all mmd files and generate svg and png files only when mmd contents change (using md5sum)
function generate_mermaid_images() {
	set -euo pipefail
	mmd_file="$1"
	# generate svg and png files only when mmd contents change (using md5sum)
	#shellcheck disable=SC2016
	md5sum=$(md5sum "$mmd_file" | $AWK '{print $1}')
	#shellcheck disable=SC2016
	md5sum_file=$(cat "$mmd_file".md5sum 2>/dev/null)
	if [ "$md5sum" != "$md5sum_file" ] || [ "$force" == true ]; then
		echo "generating svg and png files for $mmd_file"
		# generate svg and png files
		# loop formats svg and png
		# shellcheck disable=SC2043
		for fmt in svg; do #png; do
			npx -p @mermaid-js/mermaid-cli mmdc \
				--quiet --input "$mmd_file" --configFile "$docsdir"/tools/mermaid/mermaid-theme.json --cssFile "$docsdir"/tools/mermaid/mermaid.css --scale 2 --output "$mmd_file"."$fmt" --backgroundColor transparent
		done
		# update md5sum
		echo "$md5sum" >"$mmd_file".md5sum
	fi
}

export -f generate_mermaid_images

parallel -j4 --no-notice --bar --eta generate_mermaid_images ::: "$($FIND "$docsdir"/content -type f -name "*.mmd")"
