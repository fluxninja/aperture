#!/usr/bin/env bash

git_root=$(git rev-parse --show-toplevel)

FIND="find"
if [ "$(uname)" == "Darwin" ]; then
	FIND="gfind"
fi

# find all directories with a validate.sh script and save them to dirs
dirs=$($FIND "$git_root" -name validate.sh -exec dirname {} \;)

# use parallel command to cd into each directory and run validate.sh
parallel -j4 --no-notice --bar --eta "cd {} && ./validate.sh" ::: "$dirs"
