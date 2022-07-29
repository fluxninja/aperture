#!/usr/bin/env bash

# list all directories that contain .go files except for vendor directories
dirs=$(grep --include="*.go" --exclude-dir="vendor" -r "package main" -l plugins | xargs dirname | sort -u)

# use parallel to execute "cd {} && go build" in for each directory in $dirs
parallel --no-notice --bar --eta "cd {} && go build --buildmode=plugin" ::: "$dirs"
