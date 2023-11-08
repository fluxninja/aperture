#!/usr/bin/env bash
set -euo pipefail

aperturectl=$1
blueprints_uri=$2
values_file=$3

# delete the temp dir
_GEN_DIR="/tmp/aperture/_gen"
rm -rf -- "$_GEN_DIR"

"${aperturectl}" dashboard --uri "${blueprints_uri}" \
	--policy-file "${values_file}" --output-dir "${_GEN_DIR}" >&2

dashboard_dir="${_GEN_DIR}/dashboards"
# check if the dashboard dir exists
if [ ! -d "${dashboard_dir}" ]; then
	echo >&2 "Could not find dashboard directory: ${dashboard_dir}"
	exit 1
fi

# get the list of dashboards as a absolute path
dashboards=$(find "${dashboard_dir}" -type f -name '*.json' -print0 | xargs -0)

echo "${dashboards}"
