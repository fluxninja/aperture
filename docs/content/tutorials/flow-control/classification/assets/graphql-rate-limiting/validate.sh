#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$git_root"/docs/tools/aperturectl/validate_common.sh

generate_compare \
	policies/static-rate-limiting \
	values.yaml \
	tmp/policies/graphql-static-rate-limiting-cr.yaml \
	graphql-rate-limiting-jwt.yaml