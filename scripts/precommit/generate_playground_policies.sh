#!/usr/bin/env bash
set -euo pipefail
curr_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
func_dir="../"
# shellcheck source=/dev/null
source "$curr_dir/$func_dir/limit_jobs.sh"

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
		blueprint_name=$(jq --raw-output '.blueprint_name' <<<"$policy")
		values_file=$(jq --raw-output '.values_file' <<<"$policy")
		echo "Generating policies: $policy_name"
		cr_file="$scenario_dir"/policies/"$policy_name"-cr.yaml
		"$gitroot"/playground/scripts/render-policy.sh "$scenario_dir" "$aperturectl" "$gitroot/blueprints" "$blueprint_name" "$policy_name" "$scenario_dir"/"$values_file" >"$cr_file"
		"$gitroot"/scripts/git_add_safely.sh "$cr_file"
	done
}

export -f generate_policies

while IFS= read -r -d '' file
do
    limit_jobs 8 generate_policies "$file"
done < <($FIND playground -type f -name metadata.json -print0)

wait  # Wait for all background jobs to complete


popd >/dev/null
