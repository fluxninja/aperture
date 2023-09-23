#!/usr/bin/env bash
set -euo pipefail

base_dir=$1
aperturectl=$2
blueprints_uri=$3
blueprint_name=$4
values_file=$5

_GEN_DIR="${base_dir}/_gen"
trap 'rm -rf -- "$_GEN_DIR"' EXIT

"${aperturectl}" blueprints generate --name "${blueprint_name}" --uri "${blueprints_uri}" \
	--values-file "${values_file}" --output-dir "${_GEN_DIR}" >&2

dashboard_dir="${_GEN_DIR}/dashboards"
# check if the dashboard dir exists
if [ ! -d "${dashboard_dir}" ]; then
	echo >&2 "Could not find dashboard directory: ${dashboard_dir}"
	exit 1
fi

echo "${dashboard_dir}"
