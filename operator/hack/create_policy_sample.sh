#!/usr/bin/env bash
set -euo pipefail

# find git root
pushd "$(git rev-parse --show-toplevel)" >/dev/null

cp docs/content/tutorials/flow-control/concurrency-limiting/assets/basic-concurrency-limiting/basic-concurrency-limiting.yaml operator/config/samples/fluxninja.com_v1alpha1_policy.yaml

popd >/dev/null
