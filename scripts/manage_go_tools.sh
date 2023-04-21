#!/usr/bin/env bash

set -euo pipefail

git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"

install_go_tools() {
printf 'Installing Go tools\n'
go env
# install go tools
pushd "${git_root}" >/dev/null
make go-mod-tidy && make install-go-tools
popd >/dev/null
}

if ! command -v go &>/dev/null; then
  printf 'Installing Go\n'
  pushd "${git_root}" >/dev/null
  ./scripts/manage_asdf_tools.sh setup golang
  if asdf where golang &>/dev/null; then
		asdf reshim golang
	fi
  popd >/dev/null
fi

install_go_tools
