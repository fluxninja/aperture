#!/usr/bin/env bash

tools=$(grep _ ./pkg/tools/tools.go | awk -F'"' '{print $2}')

# loop $tools and call go install on each tool
for tool in $tools; do
	echo "Installing $tool"
	go install "$tool"
done
