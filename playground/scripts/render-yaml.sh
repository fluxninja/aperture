#!/usr/bin/env bash
#
# Helper script to render aperture scenarios in tilt environment
#
# Usage: render-yaml.sh path/to/scenario/manifests/

for f in "${1}"/*.yaml; do
    echo "---"
    envsubst < "$f"
done
