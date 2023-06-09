#!/usr/bin/env bash

git_root=$(git rev-parse --show-toplevel)
# Get the directory of the main script
curr_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
func_dir="../../../scripts/"

# shellcheck source=/dev/null
source "$curr_dir/$func_dir/limit_jobs.sh"

FIND="find"
if [ "$(uname)" == "Darwin" ]; then
	FIND="gfind"
fi

aperturectl="$("$git_root"/scripts/build_aperturectl.sh)"

# pull blueprints
"$aperturectl" blueprints pull --uri "$git_root"/blueprints

# find all directories with a validate.sh script and save them to dirs
dirs=$($FIND "$git_root" -name validate.sh -exec dirname {} \;)

# Use the limit_jobs function
for dir in $dirs; do
  limit_jobs 8 bash -c "cd $dir && ./validate.sh"
done

wait  # Wait for all background jobs to complete
