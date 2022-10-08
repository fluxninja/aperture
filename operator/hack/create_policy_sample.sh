#!/usr/bin/env bash
set -euo pipefail

# find git root
pushd "$(git rev-parse --show-toplevel)" >/dev/null

cp blueprints/blueprints/latency-gradient/example/gen/policies/example.yaml operator/config/samples/fluxninja.com_v1alpha1_policy.yaml

popd >/dev/null
