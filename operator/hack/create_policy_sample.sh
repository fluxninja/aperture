#!/usr/bin/env bash
set -euo pipefail

# find git root
pushd "$(git rev-parse --show-toplevel)" >/dev/null

cp docs/content/use-cases/adaptive-service-protection/assets/average-latency-feedback/policy.yaml operator/config/samples/fluxninja.com_v1alpha1_policy.yaml

popd >/dev/null
