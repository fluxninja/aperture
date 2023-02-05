#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$git_root"/docs/tools/aperturectl/validate_common.sh

generate_compare \
	"$git_root"/blueprints \
	policies/latency-aimd-concurrency-limiting \
	values.yaml \
	tmp/policies/service1-demo-app-cr.yaml \
	dynamic-rate-limiting.yaml
