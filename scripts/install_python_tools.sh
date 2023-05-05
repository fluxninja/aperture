#!/usr/bin/env bash

set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "${gitroot}" >/dev/null

if ! command -v python &>/dev/null; then
	printf 'Installing Python\n'
	./scripts/install_asdf_tools.sh setup python
fi

printf 'Installing Python tools\n'
pip3 install -r requirements.txt

if asdf where python &>/dev/null; then
	asdf reshim python
fi

popd >/dev/null
