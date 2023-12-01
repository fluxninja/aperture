#!/usr/bin/env bash

set -euo pipefail

retry_counter=10
errcode=1
while true; do
    (( retry_counter-- )) || break

    rm -rf /tmp/aperture
    git clone git@github.com:fluxninja/aperture.git /tmp/aperture
    cd /tmp/aperture
    git config user.name "FluxNinja Ops"
    git config user.email ops@fluxninja.com
    git pull
    cd sdks/aperture-go
    go get -u github.com/fluxninja/aperture/api/v2@"${COMMIT_SHA}"
    go mod tidy
    git add go.mod go.sum
    git commit -m "Update API Version in Go SDK"

    set +e
    if git push origin main; then
        errcode=0
        break
    else
        errcode="${?}"
    fi
    set -e
    sleep 1
done

exit $errcode
