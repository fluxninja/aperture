#!/usr/bin/env bash

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null || exit 1

# list all directories that contain .go files except for vendor directories
dirs=$(grep --include="*.go" --exclude-dir="vendor" -r "package main" -l plugins | xargs dirname | sort -u)

# use parallel to execute "cd {} && go build" in for each directory in $dirs
#parallel -j8 --no-notice --bar --eta "cd {} && go build --buildmode=plugin" ::: "$dirs"

echo "$dirs" | while IFS= read -r dir; do
    (cd "$dir" && go build --buildmode=plugin) &
done

wait  # Wait for all background jobs to complete



popd >/dev/null || exit 1
