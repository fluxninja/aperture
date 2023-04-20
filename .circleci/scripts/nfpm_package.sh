#!/usr/bin/env bash

set -euo pipefail
set -x

mkdir -p dist/packages/
vless_version="${APERTURECTL_BUILD_VERSION##v}"
export APERTURECTL_BUILD_VERSION=${vless_version//[-]/\~}
nfpm_args=(
	--target dist/packages/
	--config packaging/"${COMPONENT}"/nfpm.yaml
)
nfpm package "${nfpm_args[@]}" --packager deb
nfpm package "${nfpm_args[@]}" --packager rpm
ls -l dist/packages/
