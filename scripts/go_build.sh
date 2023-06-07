#!/usr/bin/env bash

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null || exit 1

# list all directories that contain .go files except for vendor directories
dirs=$(grep --include="*.go" --exclude-dir="vendor" --exclude-dir="plugins" -r "package main" -l | xargs dirname | sort -u)


echo "$dirs" | while IFS= read -r dir; do
    (cd "$dir" && go build) &
done

wait  # Wait for all background jobs to complete

popd >/dev/null || exit 1
