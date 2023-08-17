#!/usr/bin/env bash
set -euo pipefail
export LOGURU_LEVEL=TRACE
export GIT_SSH_COMMAND="fn circleci ssh -o IdentitiesOnly=yes -o IdentityAgent=none"

commit_author=$(git show --format="%aN <%aE>" --quiet)

args=(
    --author "${commit_author}"
    --release-train "${RELEASE_TRAIN:-latest}"
    --manifests-repo-url "${MANIFESTS_REPO}"
    --manifests-base-branch "${MANIFESTS_BRANCH}"
    --manifests-repo-ref "${MANIFESTS_BRANCH}"
    --push
)

if [ -n "${COMPONENT:-}" ]; then
    args+=(--component "${COMPONENT}")
fi

if [ -n "${SKIP_COMPONENT:-}" ]; then
    args+=(--skip-component "${SKIP_COMPONENT}")
fi

if [[ "${UPDATE:-}" == *","* ]]; then
    IFS=',' read -r -a updates <<< "${UPDATE}"
    for update in "${updates[@]}"; do
        args+=(--update "${update}")
    done
elif [ -n "${UPDATE:-}" ]; then
    args+=(--update "${UPDATE}")
fi

retry_counter=10
errcode=1
while true; do
    (( retry_counter-- )) || break

    set +e
    if fn --config-path "${JOB_ROOT}"/project/.opsninja.yaml \
          manifests update "${args[@]}" "${ENVIRONMENT_PATH}"; then
        errcode=0
        break
    else
        errcode="${?}"
    fi
    set -e
    sleep 1
done

exit $errcode
