#!/usr/bin/env bash
set -uo pipefail

javasdk=$(git rev-parse --show-toplevel)/sdks/aperture-java
cd "$javasdk" || exit 1

./gradlew spotlessCheck
code=$?

./gradlew spotlessApply

exit $code
