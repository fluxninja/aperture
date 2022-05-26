#!/usr/bin/env bash
set -euo pipefail

# ensure working dir - root of this repo
cd "$(dirname "$0")/../../"
pwd

# install go tools
make install-go-tools
