#!/usr/bin/env bash
set -euo pipefail

pushd ./blueprints >/dev/null
make gen-blueprints
popd >/dev/null
