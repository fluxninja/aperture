#!/usr/bin/env bash
set -euo pipefail

base_dir=$1
aperturectl=$2
blueprints_uri=$3
policy_name=$4
values_file=$5
api_key=${6:-}
endpoint=${7:-}
agent_group=${8:-default}
action=${9:-apply}
skipverify=${10:-false}
project_name=${11:-}

if [[ "${skipverify}" == "true" ]]; then
	skipverify="--skip-verify"
else
	skipverify=""
fi

SED="sed"
if [[ "$OSTYPE" == "darwin"* ]]; then
	SED="gsed"
fi

_GEN_DIR="${base_dir}/_gen"
mkdir -p "${_GEN_DIR}"
trap 'rm -rf -- "$_GEN_DIR"' EXIT

if [[ "${api_key}" != '' && "${endpoint}" != '' ]]; then
	cp "${values_file}" "${_GEN_DIR}/values.yaml"
	new_policy_name=${policy_name}-${agent_group}
	$SED -i "s/\bagent_group: .*/agent_group: ${agent_group}/g" "${_GEN_DIR}/values.yaml"
	$SED -i "s/\bpolicy_name: .*/policy_name: ${new_policy_name}/g" "${_GEN_DIR}/values.yaml"

	"${aperturectl}" blueprints generate --uri "${blueprints_uri}" --values-file "${_GEN_DIR}/values.yaml" --output-dir "${_GEN_DIR}" --overwrite >&2

	rendered_policy="${_GEN_DIR}/policies/${new_policy_name}-cr.yaml"
	if [[ "${action}" == "apply" ]]; then
		"${aperturectl}" cloud apply policy --file "${rendered_policy}" --controller "${endpoint}" --api-key "${api_key}" "${skipverify}" --project-name "${project_name}" -f -s >&2
	else
		"${aperturectl}" cloud delete policy --policy "${new_policy_name}" --controller "${endpoint}" --api-key "${api_key}" "${skipverify}" --project-name "${project_name}" || exit 0 >&2
	fi
else
	"${aperturectl}" blueprints generate --uri "${blueprints_uri}" \
		--values-file "${values_file}" --output-dir "${_GEN_DIR}" --overwrite >&2

	rendered_policy="${_GEN_DIR}/policies/${policy_name}-cr.yaml"
	if [ ! -f "${rendered_policy}" ]; then
		echo >&2 "Could not find policy file: ${rendered_policy}"
		exit 1
	fi
	head -n 1 "${rendered_policy}" | grep -q '^#' && $SED -i '1d' "${rendered_policy}"
	cat "${rendered_policy}"
fi
