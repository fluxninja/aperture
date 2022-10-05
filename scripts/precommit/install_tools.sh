#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot" >/dev/null

make install-go-tools
make install-python-tools

popd >/dev/null
