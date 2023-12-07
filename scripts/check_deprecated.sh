#!/usr/bin/env bash

GREP="grep"
if [[ "$OSTYPE" == "darwin"* ]]; then
	GREP="ggrep"
fi

# Script searches for "Deprecated: <version>" in the codebase and exits
# with an error if it finds that version is less than the latest version
# that it gets from scripts/latest_aperture_version.sh

git_root=$(git rev-parse --show-toplevel)
script_dir=$(dirname "$0")

latest_version=$("$script_dir/latest_aperture_version.sh")

ok=true

verlt() {
    [ "$1" = "$(echo -e "$1\n$2" | sort -V | head -n1)" ] && return 0 || return 1
}

# search for "Deprecated: v<version>" in the codebase for all go and proto files
# and exit with an error if it finds that version is less than the latest version
files=$($GREP -r -l "Deprecated: v" "$git_root" | grep -E "\.(go|proto)$")
# multiple lines in the file can be deprecated, so we need to check each line
for file in $files; do
	# get list of deprecated versions
	deprecated_versions=$($GREP -o "Deprecated: v[0-9]\+\.[0-9]\+\.[0-9]\+" "$file" | cut -d " " -f 2 | cut -c 2-)
	for version in $deprecated_versions; do
		if verlt "$version" "$latest_version"; then
			echo "❌ Deprecated version $version is less than latest version $latest_version in $file"
			ok=false
		else
			echo "✅ Deprecated version $version is greater than latest version $latest_version in $file"
		fi
	done
done

if [[ "$ok" == "false" ]]; then
	exit 1
fi
