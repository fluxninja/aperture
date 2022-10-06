#!/bin/bash

gitroot=$(git rev-parse --show-toplevel)
docsdir=$gitroot/docs

GREP="grep"
SED="sed"
AWK="awk"

# check whether we are using macOS
if [ "$(uname)" == "Darwin" ]; then
	GREP="ggrep"
	SED="gsed"
	AWK="gawk"
fi

# script that parses all markdown (md) files recursively in a directory
# and extracts content for the mermaid sections that begin with ```mermaid and end with ```
# and generates mmd files for each mermaid section
rm -rf tmp
mkdir -p tmp
files=$(find "$docsdir"/content -type f -name "*.md")
for f in $files; do
	dir=$(dirname "$f")
	filename=$(basename "$f")
	filenameNoExt="${filename%.*}"
	out_dir="$dir"/assets/gen/"$filenameNoExt"
	mkdir -p "$out_dir"
	rm -f "$out_dir"/*.mmd
	# extract all mermaid multiline sections that start with ```mermaid and end with ``` into an array of sections
	#shellcheck disable=SC2002,SC2016
	cat "$f" | $SED -n '/```mermaid/,/```/p' | $GREP -vP '^```$' >tmp/records.txt
	# use awk to separate out mermaid_records using RS='```' into an array of sections
	#shellcheck disable=SC2016
	$AWK '{RS="```mermaid"} NR > 1 { print $0 > "tmp/mermaid_section_" ++i}' tmp/records.txt
	# for each mermaid section, generate a mmd file
	mermaid_section_files=$(find tmp -type f -name "mermaid_section_*")
	count=0
	for mermaid_section_file in $mermaid_section_files; do
		# search for name in the comment - "%% name: <name>"
		# if found, use the name as the mmd file name
		name=$($GREP -P '^%% name: ' "$mermaid_section_file" | $SED -e 's/%% name: //')
		if [ -n "$name" ]; then
			outfilename="$name.mmd"
		else
			outfilename=$(basename "$f")_$count.mmd
		fi
		# generate mmd
		echo "generating $outfilename"
		mv "$mermaid_section_file" "$out_dir"/"$outfilename"
		git add "$out_dir"/"$outfilename"
		# increment count
		count=$((count + 1))
	done
	rm -rf tmp/*
done

rm -rf tmp
