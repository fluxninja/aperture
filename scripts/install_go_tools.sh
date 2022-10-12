#!/usr/bin/env bash

set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot"/tools/go >/dev/null

tools=$(grep _ tools.go | awk -F'"' '{print $2}')

# loop $tools and call go install on each tool
for tool in $tools; do
	echo "Installing $tool"
	go install "$tool"
done

popd >/dev/null
