#!/usr/bin/env bash

set -euo pipefail

git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"

pushd "${git_root}" >/dev/null

if ! command -v go &>/dev/null; then
	printf 'Installing Go\n'
	./scripts/manage_asdf_tools.sh setup golang
fi

printf 'Installing Go tools\n'
go env
# install go tools
make install-go-tools

if asdf where golang &>/dev/null; then
	asdf reshim golang
fi

popd >/dev/null
