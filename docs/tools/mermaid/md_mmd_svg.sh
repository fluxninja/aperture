#!/bin/bash

# script that parses all markdown (md) files recursively in a directory
# and extracts content for the mermaid sections that begin with ```mermaid and end with ```
# and generates svg files for each mermaid section
rm -f ./content/assets/gen/*.svg || true
dirs=$(find ./content -type d)
for d in $dirs; do
	files=$(find "$d" -type f -name "*.md")
	for f in $files; do
		# extract all mermaid multiline sections that start with ```mermaid and end with ``` into an array of sections
		#shellcheck disable=SC2002,SC2016
		cat "$f" | sed -n '/```mermaid/,/```/p' | grep -vP '^```$' >records.txt
		# use awk to separate out mermaid_records using RS='```' into an array of sections
		awk '{RS="```mermaid"} NR > 1 { print $0 > "mermaid_section_" ++i}' records.txt
		# for each mermaid section, generate a svg file
		mermaid_section_files=$(find . -type f -name "mermaid_section_*")
		count=0
		for mermaid_section_file in $mermaid_section_files; do
			# # extract the mermaid section content
			# mermaid_section_content=$(cat $mermaid_section_file)
			# echo "section: $mermaid_section_content"
			# generate svg
			outfilename=$(basename "$f")_$count.svg
			echo "generating $outfilename"
			npx -p @mermaid-js/mermaid-cli mmdc --configFile ./tools/mermaid/mermaid-theme.json -i "$mermaid_section_file" -o ./content/assets/gen/"$outfilename" --backgroundColor transparent
			# increment count
			count=$((count + 1))
		done
		rm -f records.txt || true
		rm -f mermaid_section_* || true
	done
done
git add ./content/assets/gen/*.svg
