#!/usr/bin/env bash

# list all directories that contain .go files except for vendor directories
dirs=$(grep --include="*.go" --exclude-dir="vendor" -r "go:generate" -l | xargs dirname | sort -u)

# use parallel to execute "cd {} && go generate" in for each directory in $dirs
parallel --no-notice --bar --eta "cd {} && go generate" ::: "$dirs"
