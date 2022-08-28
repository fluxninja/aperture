#!/usr/bin/env bash
set -euo pipefail

pushd ./api >/dev/null
make buf-generate
popd >/dev/null

pushd ./libsonnet >/dev/null
make gen-lib
popd >/dev/null
