#!/usr/bin/env bash
set -euo pipefail
TOP_LEVEL="$(git rev-parse --show-toplevel)"
# This is required to work within pre-commit `commit` stage hook
# See: https://github.com/pre-commit/pre-commit/issues/2295
# Within `commit` hook, the env variable `GIT_DIR` is set,
# which causes git dir discovery to be disabled and always return current directory
# This works because:
# 1. Outside of pre-commit, GIT_DIR is not set (so extracting TOP_LEVEL above works correctly)
# 2. Within pre-commit hook, TOP_LEVEL will be set to PWD, but we know hooks start with PWD set to git root
if [ -n "${GIT_DIR:-}" ]; then
	export GIT_WORK_TREE="${TOP_LEVEL}"
fi

# Define a minimal set of  NINJA_ variables needed to render basic development jsonnet
# manifests.
export NINJA_DEV_SETUP=true
export NINJA_KUBE_API_SERVER="localhost"

cd "${TOP_LEVEL}"/playground/tanka/
jb install
helm dependency update "${TOP_LEVEL}"/manifests/charts/aperture-agent
helm dependency update "${TOP_LEVEL}"/manifests/charts/aperture-controller
tk tool charts vendor

exit_code=0
while read -r app; do
	env_dir=$(dirname "$app")
	tk show --ext-str=CLOUD_EXTENSION="" --ext-str=VALUES="" --dangerous-allow-redirect --ext-str=projectRoot="${TOP_LEVEL}"/playground/tanka/ "$env_dir" >/dev/null || {
		exit_code="$?"
		printf '\n##########\nFAILED SHOWING ENV %s\n##########\n' "${env_dir}" >&2
	}
done < <(find apps/ -name main.jsonnet)
exit "${exit_code}"
