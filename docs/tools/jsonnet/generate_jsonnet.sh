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
	rm -rf "$out_dir"/*.jsonnet || true

	#shellcheck disable=SC2002,SC2016
	cat "$f" | $SED -n '/```jsonnet/,/```/p' | $GREP -vP '^```$' >tmp/records.txt
	# use awk to separate out jsonnet_records using RS='```' into an array of sections
	#shellcheck disable=SC2016
	$AWK '{RS="```jsonnet"} NR > 1 { print $0 > "tmp/jsonnet_section_" ++i}' tmp/records.txt
	# for each jsonnet section in tmp directory
	jsonnet_section_files=$(find tmp -type f -name "jsonnet_section_*" | sort -n)
	count=0
	for jsonnet_section_file in $jsonnet_section_files; do
		echo "Processing $f :: $jsonnet_section_file"
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

		git add "$jsonnetfilepath"

		# unset fail on error
		set +e
		# increment count
		count=$((count + 1))
	done
	rm -rf tmp/*
done

# find all jsonnet files in docs/content directory and generate them
jsonnet_files=$(find "$docsdir"/content -type f -name "*.jsonnet")
for jsonnet_file in $jsonnet_files; do
	echo "Processing $jsonnet_file"
	dir=$(dirname "$jsonnet_file")
	# remove extension and add .yaml
	yamlfilepath="${jsonnet_file%.*}".yaml

	# cp jsonnet file to tmp file
	tmpjsonnetfilepath=tmp/"$(basename "$jsonnet_file")"
	cp "$jsonnet_file" "$tmpjsonnetfilepath"

	old_yaml_file_contents=""
	if [ -f "$yamlfilepath" ]; then
		old_yaml_file_contents=$(cat "$yamlfilepath")
	fi

	# replace github.com/fluxninja/aperture/blueprints with $"gitroot"/blueprints
	$SED -i "s|github.com/fluxninja/aperture/blueprints|$gitroot/blueprints|g" "$tmpjsonnetfilepath"
	# fail script if any of the below commands fail
	set -e
	jsonnet -J "$gitroot"/blueprints/vendor "$tmpjsonnetfilepath" | yq -P >"$yamlfilepath"
	# run prettier
	npx prettier --write "$yamlfilepath"
	git add "$yamlfilepath"

	# if the file is a policy kind then generate mermaid diagram
	if [ "$(yq e '.kind == "Policy"' "$yamlfilepath")" = "true" ]; then
		# generate mermaid diagram
		# compile the policy and generate mermaid if yaml has changed
		if [ "$old_yaml_file_contents" != "$(cat "$yamlfilepath")" ]; then
			# generate mermaid diagram
			mermaidfilepath="${jsonnet_file%.*}".mmd
			# compile the policy
			go run "$gitroot"/cmd/circuit-compiler/main.go -cr "$yamlfilepath" --mermaid "$mermaidfilepath"
			git add "$mermaidfilepath"
		fi
	fi
	# unset fail on error
	set +e
	rm -rf tmp/*
done

rm -rf tmp
