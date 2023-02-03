#!/usr/bin/env bash

set -euo pipefail
set -x

mkdir -p dist/packages/
vless_version="${VERSION##v}"
export VERSION=${vless_version//[-]/\~}
nfpm_args=(
  --target dist/packages/
  --config packaging/"${COMPONENT}"/nfpm.yaml
)
nfpm package "${nfpm_args[@]}" --packager deb
nfpm package "${nfpm_args[@]}" --packager rpm
ls -l dist/packages/
