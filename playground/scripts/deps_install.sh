#!/usr/bin/env bash
#
# A simple wrapper over `jb install` that installs libsonnet and blueprints under
# expected paths. This is a workaround for https://github.com/jsonnet-bundler/jsonnet-bundler/issues/70
set -euo pipefail

BASE_DIR=$1

(
cd "${BASE_DIR}"
jb install
)

GIT_REPO_ROOT=$(git rev-parse --show-toplevel)

VENDOR_PATH=${BASE_DIR}/vendor/github.com/fluxninja/aperture/

symlink_directory() {
    local directory=$1

    if [ ! -d "${VENDOR_PATH}" ]; then
        mkdir -p "${VENDOR_PATH}"
    fi

    if [[ ! -L "${VENDOR_PATH}"/"${directory}" ]] || [[ ! -e "${VENDOR_PATH}"/"${directory}" ]]; then
        [ -e "${VENDOR_PATH}"/"${directory}" ] && rm -rf "${VENDOR_PATH:?}"/"${directory}"
        echo "${GIT_REPO_ROOT}/${directory} -> ${VENDOR_PATH}/${directory}"
        ln -s "${GIT_REPO_ROOT}"/"${directory}" "${VENDOR_PATH}"/"${directory}"
    fi
}

symlink_directory blueprints

if [ -f "${BASE_DIR}/chartfile.yaml" ]; then
    pushd "${BASE_DIR}" > /dev/null
    tk tool charts vendor
    popd > /dev/null
fi


CHARTS_DIR="${BASE_DIR}/charts"
if [ ! -d "${CHARTS_DIR}" ]; then
    exit 0
fi

for path in "${CHARTS_DIR}"/*; do
    # Process directories and directory symlinks, skipping others
    if [ ! -d "$path" ]; then
        continue
    fi
    pushd "$path" > /dev/null
    if helm dependency list | grep -q missing$; then
        helm dependency build
    fi
    popd > /dev/null
done
