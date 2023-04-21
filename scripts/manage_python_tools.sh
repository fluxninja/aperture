#!/usr/bin/env bash

set -euo pipefail

git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"

install_python_tools() {
printf 'Installing Python tools\n'
pushd "${git_root}" >/dev/null
make install-python-tools
popd >/dev/null
}

if ! command -v python &>/dev/null; then
  printf 'Installing Python\n'
  pushd "${git_root}" >/dev/null
  ./scripts/manage_asdf_tools.sh setup python
  popd >/dev/null

  if asdf where python &>/dev/null; then
		asdf reshim python
	fi
fi

install_python_tools
