#!/usr/bin/env bash

set -euo pipefail

git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"

pushd "${git_root}" >/dev/null
./scripts/manage_asdf_tools.sh setup
./scripts/manage_go_tools.sh
./scripts/manage_python_tools.sh
popd >/dev/null
