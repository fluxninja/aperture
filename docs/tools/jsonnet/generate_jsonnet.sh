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

function generate_jsonnet_files() {
	set -euo pipefail
	jsonnet_file="$1"
	echo "Processing $jsonnet_file"
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
