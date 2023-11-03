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

# Prepare list of commands to run in parallel
declare -a cmds=()
for dir in $dirs; do
	cmds+=("cd $dir && ./validate.sh")
done

# Run commands in parallel
"$git_root"/scripts/run_parallel.sh "${cmds[@]}"
