#!/usr/bin/env bash
set -euo pipefail

blueprint_name=${1}
echo "$blueprint_name"

script_root=$(dirname "$0")
blueprints_root=${script_root}/..

python "${blueprints_root}"/scripts/blueprint-readme-generator.py "${blueprints_root}"/blueprints/"${blueprint_name}"
if which prettier > /dev/null 2>&1; then
    prettier --write "${blueprints_root}"/blueprints/"${blueprint_name}"/README.md
fi

python "${blueprints_root}"/scripts/aperture-generate.py --output "${blueprints_root}"/_gen/ \
    --config "${blueprints_root}"/examples/demoapp-"${blueprint_name}".jsonnet \
    "${blueprints_root}"/blueprints/"${blueprint_name}"

go run -mod=mod "${blueprints_root}"/../cmd/circuit-compiler/main.go \
    -policy "${blueprints_root}"/_gen/policies/service1-"${blueprint_name}".yaml \
    -dot "${blueprints_root}"/blueprints/"${blueprint_name}"/graph.dot
dot -Tsvg "${blueprints_root}"/blueprints/"${blueprint_name}"/graph.dot \
     > "${blueprints_root}"/blueprints/"${blueprint_name}"/graph.svg
