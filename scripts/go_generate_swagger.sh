#!/usr/bin/env bash

# check whether ggrep exists or not
if ! [ -x "$(command -v ggrep)" ]; then
	GREP="grep"
else
	GREP="ggrep"
fi

# check whether gdirname exists or not
if ! [ -x "$(command -v gdirname)" ]; then
	DIRNAME="dirname"
else
	DIRNAME="gdirname"
fi

# list all directories that contain .go files except for vendor directories
dirs=$("${GREP}" --include="*.go" --exclude-dir="vendor" -r "go:generate swagger" -l | xargs "${DIRNAME}" | sort -u)

# use parallel to execute "cd {} && go generate" in for each directory in $dirs
parallel --no-notice --bar --eta "cd {} && go generate" ::: "$dirs"
