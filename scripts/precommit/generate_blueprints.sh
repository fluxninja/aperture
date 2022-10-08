#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/blueprints >/dev/null
make generate-blueprints
popd >/dev/null
