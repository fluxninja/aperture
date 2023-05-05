#!/usr/bin/env bash

set -euo pipefail

git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"

pushd "${git_root}" >/dev/null

if ! command -v python &>/dev/null; then
	printf 'Installing Python\n'
	./scripts/manage_asdf_tools.sh setup python
fi

printf 'Installing Python tools\n'
make install-python-tools

if asdf where python &>/dev/null; then
	asdf reshim python
fi

popd >/dev/null
