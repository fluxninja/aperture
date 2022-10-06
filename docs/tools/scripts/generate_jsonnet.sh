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

rm -rf tmp || true
mkdir -p tmp
files=$(find "$docsdir"/content -type f -name "*.md")
for f in $files; do
	dir=$(dirname "$f")
	filename=$(basename "$f")
	filenameNoExt="${filename%.*}"
	out_dir="$dir"/assets/gen/"$filenameNoExt"/jsonnet
	mkdir -p "$out_dir"
	rm -f "$out_dir"/*

	#shellcheck disable=SC2002,SC2016
	cat "$f" | $SED -n '/```jsonnet/,/```/p' | $GREP -vP '^```$' >tmp/records.txt
	# use awk to separate out jsonnet_records using RS='```' into an array of sections
	#shellcheck disable=SC2016
	$AWK '{RS="```jsonnet"} NR > 1 { print $0 > "tmp/jsonnet_section_" ++i}' tmp/records.txt
	# for each jsonnet section in tmp directory
	jsonnet_section_files=$(find tmp -type f -name "jsonnet_section_*")
	count=0
	for jsonnet_section_file in $jsonnet_section_files; do
		echo "Processing $f :: $jsonnet_section_file"
		outfilename="$out_dir"/"$filenameNoExt"_"$count".yaml
		# replace github.com/fluxninja/aperture/blueprints with $"gitroot"/blueprints
		$SED -i "s|github.com/fluxninja/aperture/blueprints|$gitroot/blueprints|g" "$jsonnet_section_file"
		# fail script if any of the below commands fail
		set -e
		jsonnet --yaml-stream -J "$gitroot"/blueprints/vendor "$jsonnet_section_file" >"$outfilename"
		if [ "$(yq e '.kind == "Policy"' "$outfilename")" = "true" ]; then
			specfilename="$out_dir"/"$filenameNoExt"_"$count"_spec.yaml
			mermaidfilename="$out_dir"/"$filenameNoExt"_"$count".mmd
			# extract spec key from yaml
			yq '.spec' "$outfilename" >"$specfilename"
			git add "$specfilename"
			# validate with circuit compiler
			go run "$gitroot"/cmd/circuit-compiler/main.go -policy "$specfilename" --mermaid "$mermaidfilename"
			git add "$mermaidfilename"

		fi
		git add "$outfilename"
		# unset fail on error
		set +e
		# increment count
		count=$((count + 1))
	done
	rm -rf tmp/*
done

rm -rf tmp
