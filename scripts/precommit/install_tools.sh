#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot" >/dev/null

make install-go-tools
make install-python-tools

if asdf current golang >/dev/null 2>/dev/null; then
  asdf reshim golang
fi

popd >/dev/null
