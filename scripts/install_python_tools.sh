#!/usr/bin/env bash

set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "${gitroot}" >/dev/null

if ! command -v python &>/dev/null; then
	printf 'Installing Python\n'
	./scripts/install_asdf_tools.sh setup python
fi

printf 'Installing Python tools\n'
# remove once https://github.com/yaml/pyyaml/issues/601 is fixed
pip uninstall -y pyyaml
pip install "Cython<3.0" PyYAML==5.4.1 --no-build-isolation
pip3 install -r requirements.txt

if asdf where python &>/dev/null; then
	asdf reshim python
fi

popd >/dev/null
