#!/usr/bin/env bash
set -euo pipefail

# find git root
pushd "$(git rev-parse --show-toplevel)" >/dev/null

cp docs/content/guides/service-load-management/assets/average-latency-feedback/policy.yaml operator/config/samples/fluxninja.com_v1alpha1_policy.yaml

popd >/dev/null
