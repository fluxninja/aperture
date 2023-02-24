#!/usr/bin/env bash
set -euo pipefail

FIND="find"
if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi
export FIND

gitroot="$(git rev-parse --show-toplevel)"
export gitroot

aperturectl="$("$gitroot"/scripts/build_aperturectl.sh)"
export aperturectl

pushd "$gitroot" >/dev/null

function generate_policies() {
	set -euo pipefail
	scenario="$1"
	scenario_dir=$(dirname "$scenario")
	rm -rf "$scenario_dir"/**/*-cr.yaml
	metadata_file="$scenario_dir"/metadata.json
	readarray -t policies < <(jq --compact-output '.aperture_policies[]' "$metadata_file")
	for policy in "${policies[@]}"; do
		policy_name=$(jq --raw-output '.policy_name' <<<"$policy")
		policy_def=$(jq --raw-output '.policy_def' <<<"$policy")
		values_file=$(jq --raw-output '.values_file' <<<"$policy")
		echo "Generating policies: $policy_name"
		"$gitroot"/playground/scripts/render-policy.sh "$scenario_dir" "$aperturectl" "$gitroot/blueprints" "$policy_def" "$policy_name" "$scenario_dir"/"$values_file" >/dev/null
	done
}

export -f generate_policies

parallel -j4 --no-notice --bar --eta generate_policies ::: "$($FIND playground -type f -name metadata.json)"

popd >/dev/null
