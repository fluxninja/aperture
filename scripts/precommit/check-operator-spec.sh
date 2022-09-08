#!/usr/bin/env bash
set -euo pipefail

make operator-generate
make operator-manifests
