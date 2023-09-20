#!/usr/bin/env bash
set -euo pipefail
set -x

base_path="https://github.com/goreleaser/nfpm/releases"
version="2.32.0"
filename="nfpm_${version}_amd64.deb"
file="/tmp/${filename}"
curl --silent --show-error --location "${base_path}/download/v${version}/${filename}" -o "${file}"
sudo dpkg -i "${file}"
