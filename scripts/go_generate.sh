#!/usr/bin/env bash

# check whether ggrep exists or not
if ! [ -x "$(command -v ggrep)" ]; then
	echo "picking grep"
	GREP="grep"
else
	echo "picking ggrep"
	GREP="ggrep"
fi

# check whether gdirname exists or not
if ! [ -x "$(command -v gdirname)" ]; then
	echo "picking dirname"
	DIRNAME="dirname"
else
	echo "picking gdirname"
	DIRNAME="gdirname"
fi

# list all directories that contain .go files except for vendor directories
dirs=$("${GREP}" --include="*.go" --exclude-dir="vendor" -r "go:generate" -l | xargs "${DIRNAME}" | sort -u)

# use parallel to execute "cd {} && go generate" in for each directory in $dirs
parallel -j4 "cd {} && go generate" ::: "$dirs"
