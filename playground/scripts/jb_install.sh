#!/usr/bin/env bash
#
# A simple wrapper over `jb install` that installs libsonnet and blueprints under
# expected paths. This is a workaround for https://github.com/jsonnet-bundler/jsonnet-bundler/issues/70
set -euo pipefail

(
cd tanka
jb install
)

GIT_REPO_ROOT=$(git rev-parse --show-toplevel)

SCRIPT_ROOT=$(dirname "$0")
VENDOR_PATH=${SCRIPT_ROOT}/../tanka/vendor/github.com/fluxninja/aperture

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

symlink_directory libsonnet
symlink_directory blueprints
