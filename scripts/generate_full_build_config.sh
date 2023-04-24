#!/usr/bin/env bash
set -eu
set -o pipefail

EXTENSIONS_DIR=${1:-$(git rev-parse --show-toplevel)/extensions}

pushd "$EXTENSIONS_DIR" > /dev/null

echo "bundled_extensions:"
find integrations/ -mindepth 2 -maxdepth 2 -type d | sort | awk '{ print "- "$1 }'

# Note: not listing fluxninja nor sentry, as they're enabled by default anyway
