#!/usr/bin/env bash
set -euo pipefail

base_dir=$1
aperturectl=$2
blueprints_uri=$3
policy_def=$4
policy_name=$5
values_file=$6

SED="sed"
if [[ "$OSTYPE" == "darwin"* ]]; then
	SED="gsed"
fi

_GEN_DIR="${base_dir}/_gen"
trap 'rm -rf -- "$_GEN_DIR"' EXIT

"${aperturectl}" blueprints generate --name "${policy_def}" --uri "${blueprints_uri}" \
	--values-file "${values_file}" --output-dir "${_GEN_DIR}" >&2

rendered_policy="${_GEN_DIR}/policies/${policy_name}-cr.yaml"
if [ ! -f "${rendered_policy}" ]; then
	echo >&2 "Could not find policy file: ${rendered_policy}"
	exit 1
fi

head -n 1 "${rendered_policy}" | grep -q '^#' && $SED -i '1d' "${rendered_policy}"
cp "${rendered_policy}" "$(dirname "${values_file}")"
cat "${rendered_policy}"
