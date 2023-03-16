#!/usr/bin/env bash
#
# This script runs go mod tidy in all directories that contain go.mod files.

FIND="find"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null || exit 1

# list all directories that contain go.mod files except for vendor directories

dirs=$($FIND . -name go.mod -not -path "*/vendor/*" -exec dirname {} \; | sort -u)

parallel -j4 --no-notice --bar --eta "cd {} && go mod tidy -compat=1.20" ::: "$dirs"

popd >/dev/null || exit 1
