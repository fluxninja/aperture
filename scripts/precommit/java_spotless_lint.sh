#!/usr/bin/env bash
set -uo pipefail

javasdk=$(git rev-parse --show-toplevel)/sdks/aperture-java
cd "$javasdk" || exit 1

code=0
if ./gradlew spotlessCheck; then
  echo "Java code is formatted correctly"
else
  code=$?
  echo "Formatting Java code"
  ./gradlew spotlessApply
fi

exit $code
