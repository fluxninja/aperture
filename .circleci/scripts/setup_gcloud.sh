#!/usr/bin/env bash
set -euo pipefail
set -x

# Setup gcloud access
echo "${GCLOUD_AR_SERVICE_KEY?}" | gcloud auth activate-service-account --key-file=-
gcloud --quiet config set project "${GOOGLE_PROJECT_ID?}"
gcloud --quiet config set compute/zone "${GOOGLE_COMPUTE_ZONE?}"
