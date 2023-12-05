#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/sdks/aperture-py >/dev/null
pydoctor
popd >/dev/null
