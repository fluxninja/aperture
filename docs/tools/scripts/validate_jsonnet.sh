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
		# replace github.com/fluxninja/aperture/blueprints with $"gitroot"/blueprints
		$SED -i "s|github.com/fluxninja/aperture/blueprints|$gitroot/blueprints|g" "$jsonnet_section_file"
		# fail script if any of the below commands fail
		set -e
		jsonnet --yaml-stream -J "$gitroot"/blueprints/vendor "$jsonnet_section_file" >tmp/output.yaml
		# check whether output.yaml contains the key "kind: Policy" i.e. output of yq is true
		if [ "$(yq e '.kind == "Policy"' tmp/output.yaml)" = "true" ]; then
			# extract spec key from yaml
			yq '.spec' tmp/output.yaml >tmp/spec.yaml
			# validate with circuit compiler
			go run "$gitroot"/cmd/circuit-compiler/main.go -policy tmp/spec.yaml
		fi
		# unset fail on error
		set +e
		# increment count
		count=$((count + 1))
	done
	rm -rf tmp/*
done

rm -rf tmp
