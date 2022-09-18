#!/bin/bash

# script that parses all markdown (md) files recursively in a directory
# and extracts content for the mermaid sections that begin with ```mermaid and end with ```
# and generates mmd files for each mermaid section
out_dir=./content/assets/gen
rm -f "$out_dir"/*.mmd || true
dirs=$(find ./content -type d)
for d in $dirs; do
	files=$(find "$d" -type f -name "*.md")
	for f in $files; do
		# extract all mermaid multiline sections that start with ```mermaid and end with ``` into an array of sections
		#shellcheck disable=SC2002,SC2016
		cat "$f" | sed -n '/```mermaid/,/```/p' | grep -vP '^```$' >records.txt
		# use awk to separate out mermaid_records using RS='```' into an array of sections
		awk '{RS="```mermaid"} NR > 1 { print $0 > "mermaid_section_" ++i}' records.txt
		# for each mermaid section, generate a mmd file
		mermaid_section_files=$(find . -type f -name "mermaid_section_*")
		count=0
		for mermaid_section_file in $mermaid_section_files; do
			# search for name in the comment - "%% name: <name>"
			# if found, use the name as the mmd file name
			name=$(grep -P '^%% name: ' "$mermaid_section_file" | sed -e 's/%% name: //')
			if [ -n "$name" ]; then
				outfilename="$name.mmd"
			else
				outfilename=$(basename "$f")_$count.mmd
			fi
			# generate mmd
			echo "generating $outfilename"
			mv "$mermaid_section_file" "$out_dir"/"$outfilename"
			# increment count
			count=$((count + 1))
		done
		rm -f records.txt || true
		rm -f mermaid_section_* || true
	done
done
git add "$out_dir"/*.mmd
