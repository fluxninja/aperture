#!/usr/bin/env bash
set -euo pipefail

pushd ./blueprints >/dev/null
make generate-blueprints
popd >/dev/null
