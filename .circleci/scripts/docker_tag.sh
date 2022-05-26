#!/usr/bin/env bash

SOURCE_TAG=$(eval echo "${PARAM_SOURCE_TAG}")
DEST_TAG=$(eval echo "${PARAM_DEST_TAG}")

source_image="${PARAM_SOURCE_REGISTRY}/${PARAM_SOURCE_IMAGE}:${SOURCE_TAG}"
dest_image="${PARAM_DEST_REGISTRY}/${PARAM_DEST_IMAGE}:${DEST_TAG}"

echo "Tagging ${source_image} -> ${dest_image}"

docker image ls
docker tag "${source_image}" "${dest_image}"
