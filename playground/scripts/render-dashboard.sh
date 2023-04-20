#!/usr/bin/env bash
set -euo pipefail

base_dir=$1
aperturectl=$2
blueprints_uri=$3
blueprint_name=$4
policy_name=$5
values_file=$6

_GEN_DIR="${base_dir}/_gen"
trap 'rm -rf -- "$_GEN_DIR"' EXIT

"${aperturectl}" blueprints generate --name "${blueprint_name}" --uri "${blueprints_uri}" \
	--values-file "${values_file}" --output-dir "${_GEN_DIR}" >&2

rendered_dashboard="${_GEN_DIR}/dashboards/${policy_name}.json"
if [ ! -f "${rendered_dashboard}" ]; then
	echo >&2 "Could not find dashboard file: ${rendered_dashboard}"
	exit 1
fi

tr -d '\n' < "${rendered_dashboard}"
