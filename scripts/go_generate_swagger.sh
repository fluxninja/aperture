#!/usr/bin/env bash

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null || exit 1

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
parallel -j8 --no-notice --bar --eta "cd {} && go generate -v -x" ::: "$dirs"

popd >/dev/null || exit 1
