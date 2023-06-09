#!/bin/bash

# Get the directory of the main script
curr_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
func_dir="../../../scripts/"
# shellcheck source=/dev/null
source "$curr_dir/$func_dir/limit_jobs.sh"

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

aperturectl="$("$gitroot"/scripts/build_aperturectl.sh)"
export aperturectl

# accept a --force flag to force the regeneration of all graphs
# check if $1 is bounded and equal to --force
if [[ -n "${1:-}" ]] && [[ "$1" == "--force" ]]; then
	force=true
else
	force=false
fi
export force

# build json2yaml
go build -o "$scriptroot"/json2yaml "$scriptroot"/json2yaml.go

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
	if [ "$(yq e '.kind == "Policy"' "$jsonfilepath" --output-format=yaml)" = "true" ]; then
		# convert the policy to yaml
		"$scriptroot"/json2yaml "$jsonfilepath" "$yamlfilepath"
		rm -rf "$jsonfilepath"
		# run prettier
		prettier --write "$yamlfilepath"
		# generate mermaid diagram
		mermaidfilepath="${jsonnet_file%.*}".mmd
		# compile the policy
		"$aperturectl" compile --cr "$yamlfilepath" --mermaid "$mermaidfilepath"
	else
		prettier --write "$jsonfilepath"
	fi
	rm -rf "$tmpdir"
}

export -f generate_jsonnet_files

# find all jsonnet files and run generate_jsonnet_files in parallel
while IFS= read -r -d '' file
do
    limit_jobs 8 generate_jsonnet_files "$file"
done < <($FIND "$docsdir"/content -type f -name '*.jsonnet' -print0)

wait  # Wait for all background jobs to complete
