#!/usr/bin/env bash

git_root=$(git rev-parse --show-toplevel)

FIND="find"
if [ "$(uname)" == "Darwin" ]; then
	FIND="gfind"
fi

aperturectl="$("$git_root"/scripts/build_aperturectl.sh)"

# pull blueprints
"$aperturectl" blueprints pull --uri "$git_root"/blueprints

# find all directories with a validate.sh script and save them to dirs
dirs=$($FIND "$git_root" -name validate.sh -exec dirname {} \;)

# use parallel command to cd into each directory and run validate.sh
parallel -j4 --no-notice --bar --eta --halt-on-error now,fail,1 "cd {} && ./validate.sh" ::: "$dirs"
