#!/usr/bin/env bash
set -euo pipefail

base_dir=$1
aperturectl=$2
blueprints_uri=$3
blueprint_name=$4
policy_name=$5
values_file=$6
api_key=${7:-}
endpoint=${8:-}
agent_group=${9:-default}
action=${10:-apply}

SED="sed"
if [[ "$OSTYPE" == "darwin"* ]]; then
	SED="gsed"
fi

_GEN_DIR="${base_dir}/_gen"
mkdir -p "${_GEN_DIR}"
trap 'rm -rf -- "$_GEN_DIR"' EXIT

if [[ "$api_key" != '' && "$endpoint" != '' ]]; then
	cp "${values_file}" "${_GEN_DIR}/values.yaml"
	new_policy_name=${policy_name}-${agent_group}
	$SED -i "s/\bagent_group: .*/agent_group: ${agent_group}/g" "${_GEN_DIR}/values.yaml"
	$SED -i "s/\bpolicy_name: .*/policy_name: ${new_policy_name}/g" "${_GEN_DIR}/values.yaml"

	"${aperturectl}" blueprints generate --name "${blueprint_name}" --uri "${blueprints_uri}" \
		--values-file "${_GEN_DIR}/values.yaml" --output-dir "${_GEN_DIR}" --overwrite >&2

	rendered_policy="${_GEN_DIR}/policies/${new_policy_name}-cr.yaml"
	if [[ "${action}" == "apply" ]]; then
		"${aperturectl}" apply policy --file "${rendered_policy}" --controller "${endpoint}" --api-key "${api_key}" -f -s >&2
	else
		"${aperturectl}" delete policy --policy "${new_policy_name}" --controller "${endpoint}" --api-key "${api_key}" >&2
	fi
else
	"${aperturectl}" blueprints generate --name "${blueprint_name}" --uri "${blueprints_uri}" \
		--values-file "${values_file}" --output-dir "${_GEN_DIR}" --overwrite >&2

	rendered_policy="${_GEN_DIR}/policies/${policy_name}-cr.yaml"
	if [ ! -f "${rendered_policy}" ]; then
		echo >&2 "Could not find policy file: ${rendered_policy}"
		exit 1
	fi
	head -n 1 "${rendered_policy}" | grep -q '^#' && $SED -i '1d' "${rendered_policy}"
	cat "${rendered_policy}"
fi
