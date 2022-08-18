#!/usr/bin/env bash
set -euo pipefail

# generate docs
# Note: this is not using generate-docs target, because of this bug:
# https://github.com/fluxninja/aperture/issues/126
# Regenerating svgs needs to be done manually via `make generate-mermaid`.
make generate-config-markdown
