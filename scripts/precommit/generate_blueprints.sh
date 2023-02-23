#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

"$gitroot"/scripts/generate_blueprints.sh
