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
		# skip this file if it contains "@include:"
		if $GREP -q "@include:" "$mermaid_section_file"; then
			continue
		fi
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

# find all mmd files and generate svg and png files only when mmd contents change (using md5sum)

# find all mmd files
mmd_files=$(find "$docsdir"/content -type f -name "*.mmd")
for mmd_file in $mmd_files; do
	# generate svg and png files only when mmd contents change (using md5sum)
	#shellcheck disable=SC2016
	md5sum=$(md5sum "$mmd_file" | $AWK '{print $1}')
	#shellcheck disable=SC2016
	md5sum_file=$(cat "$mmd_file".md5sum 2>/dev/null)
	if [ "$md5sum" != "$md5sum_file" ]; then
		echo "generating svg and png files for $mmd_file"
		# generate svg and png files
		# loop formats svg and png
		for fmt in svg png; do
			npx -p @mermaid-js/mermaid-cli mmdc \
				--quiet --input "$mmd_file" --configFile "$docsdir"/tools/mermaid/mermaid-theme.json --cssFile ./tools/mermaid/mermaid.css --scale 2 --output "$mmd_file"."$fmt" --backgroundColor transparent
			git add "$mmd_file"."$fmt"
		done
		# update md5sum
		echo "$md5sum" >"$mmd_file".md5sum
		git add "$mmd_file".md5sum
	fi
done
