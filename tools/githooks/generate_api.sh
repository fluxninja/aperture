#!/usr/bin/env bash
set -euo pipefail

cd ./api && make buf-generate
