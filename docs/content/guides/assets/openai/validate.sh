#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$git_root"/docs/tools/aperturectl/validate_common.sh

# Array of values files and their corresponding policy and graph names
declare -a values_files=("values1" "values2")
declare -a policy_names=("gpt-4-tpm-cr" "gpt-4-rpm-cr")

for i in "${!values_files[@]}"; do
    values_file=${values_files[$i]}
    policy_name=${policy_names[$i]}

    generate_from_values \
        "${values_file}.yaml" \
        tmp

    # Copy the generated policy and graph to this (assets) directory so that they can be used in the docs
    cp "tmp/policies/${policy_name}.yaml" policy.yaml
    cp "tmp/graphs/${policy_name}.mmd" graph.mmd

    # git add the generated policy and graph
    "$git_root"/scripts/git_add_safely.sh policy.yaml graph.mmd

    # Remove the tmp directory
    rm -rf tmp
done
