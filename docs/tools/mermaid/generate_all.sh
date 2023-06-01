#!/bin/bash

gitroot=$(git rev-parse --show-toplevel)
export gitroot

docsdir=$gitroot/docs

scriptroot=$(dirname "$0")
export scriptroot

FIND="find"
AWK="awk"
# check whether we are using macOS
if [ "$(uname)" == "Darwin" ]; then
	FIND="gfind"
	AWK="gawk"
fi
export FIND
export AWK

# accept a --force flag to force generation of all mermaid files
if [[ -n "${1:-}" ]] && [[ "$1" == "--force" ]]; then
	force=true
else
	force=false
fi
export force

# find all mmd files and generate svg and png files only when mmd contents change (using md5sum)
function generate_mermaid_images() {
	set -euo pipefail
	mmd_file="$1"
	# generate svg and png files only when mmd contents change (using md5sum)
	#shellcheck disable=SC2016
	md5sum=$(md5sum "$mmd_file" | $AWK '{print $1}')
	# check if md5sum file exists
	md5sum_file=""
	if [ -f "$mmd_file".md5sum ]; then
		#shellcheck disable=SC2016
		md5sum_file=$(cat "$mmd_file".md5sum 2>/dev/null)
	fi
	if [ "$md5sum" != "$md5sum_file" ] || [ "$force" == true ]; then
		echo "generating svg and png files for $mmd_file"
		"$scriptroot"/generate_mermaid.sh "$mmd_file"
		# update md5sum
		echo "$md5sum" >"$mmd_file".md5sum
	fi
}

export -f generate_mermaid_images

parallel -j8 --halt-on-error now,fail,1 --no-notice --bar --eta generate_mermaid_images ::: "$($FIND "$docsdir"/content -type f -name "*.mmd")"
