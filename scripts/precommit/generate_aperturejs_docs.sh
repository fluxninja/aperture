#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

pushd "$gitroot"/sdks/aperture-js >/dev/null
npm install
npm run docs
popd >/dev/null
