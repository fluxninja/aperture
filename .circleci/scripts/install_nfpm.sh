#!/usr/bin/env bash
set -euo pipefail
set -x

filename="nfpm_amd64.deb"
file="/tmp/${filename}"
curl --silent --show-error --location https://github.com/goreleaser/nfpm/releases/latest/download/${filename} -o "${file}"
sudo dpkg -i "${file}"
