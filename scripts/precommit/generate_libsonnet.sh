#!/usr/bin/env bash
set -euo pipefail

pushd ./libsonnet >/dev/null
make gen-lib
popd >/dev/null
