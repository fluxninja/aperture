#!/usr/bin/env bash
set -euo pipefail
set -x

base_path="https://github.com/goreleaser/nfpm/releases"
version=$(curl --silent -I "${base_path}/latest" | grep location: | awk -F '/' '{print $NF}' | tr -d 'v\r\n')
filename="nfpm_${version}_amd64.deb"
file="/tmp/${filename}"
curl --silent --show-error --location "${base_path}/download/v${version}/${filename}" -o "${file}"
sudo dpkg -i "${file}"
