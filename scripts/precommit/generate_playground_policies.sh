#!/usr/bin/env bash
set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"
aperturectl="$("$gitroot"/playground/scripts/build_aperturectl.sh)"

pushd "$gitroot" >/dev/null
find playground -type f -name metadata.json | while IFS= read -r scenario; do
    scenario_dir=$(dirname "$scenario")
    rm -rf "$scenario_dir"/**/*-cr.yaml
    metadata_file="$scenario_dir"/metadata.json
    readarray -t policies < <(jq --compact-output '.aperture_policies[]' "$metadata_file")
    for policy in "${policies[@]}"; do
        policy_name=$(jq --raw-output '.policy_name' <<< "$policy")
        policy_def=$(jq --raw-output '.policy_def' <<< "$policy")
        values_file=$(jq --raw-output '.values_file' <<< "$policy")
        echo "Generating policies: $policy_name"
        "$gitroot"/playground/scripts/render-policy.sh "$gitroot/playground" "$aperturectl" "$gitroot/blueprints" "$policy_def" "$policy_name" "$scenario_dir"/"$values_file" >/dev/null
    done
done
popd >/dev/null
