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
		blueprint_name=$(jq --raw-output '.blueprint_name' <<<"$policy")
		values_file=$(jq --raw-output '.values_file' <<<"$policy")
		echo "Generating policies: $policy_name"
		cr_file="$scenario_dir"/policies/"$policy_name"-cr.yaml
		"$gitroot"/playground/scripts/render-policy.sh "$scenario_dir" "$aperturectl" "$gitroot/blueprints" "$blueprint_name" "$policy_name" "$scenario_dir"/"$values_file" >"$cr_file"
		"$gitroot"/scripts/git_add_safely.sh "$cr_file"
	done
}

export -f generate_policies

declare -a cmds=()
while IFS= read -r -d '' file; do
	cmds+=("generate_policies '$file'")
done < <($FIND playground -type f -name metadata.json -print0)

# Run the policy generation commands in parallel
"$gitroot"/scripts/run_parallel.sh "${cmds[@]}"

popd >/dev/null
