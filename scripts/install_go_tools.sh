#!/usr/bin/env bash

gitroot=$(git rev-parse --show-toplevel)

tools=$(grep _ "$gitroot"/pkg/tools/tools.go | awk -F'"' '{print $2}')

# loop $tools and call go install on each tool
for tool in $tools; do
	echo "Installing $tool"
	go install "$tool"
done
