#!/usr/bin/env bash
set -e

DOCKER_TAGS_ARG=""

parse_tags_to_docker_arg() {
	# Set comma as the new delimiter for the scope of this function.
	local IFS=","

	# Split tags into an array based on IFS delimiter.
	read -ra tags <<<"$PARAM_TAG"

	local docker_arg

	for tag in "${tags[@]}"; do
		if [ -z "$docker_arg" ]; then
			docker_arg="--tag=\"$PARAM_REGISTRY/$PARAM_IMAGE:$tag\""
		else
			docker_arg="$docker_arg --tag=\"$PARAM_REGISTRY/$PARAM_IMAGE:$tag\""
		fi
	done

	DOCKER_TAGS_ARG=$(eval echo "${docker_arg}")
}

pull_images_from_cache() {
	local cache
	cache=$(eval echo "${PARAM_CACHE_FROM}")

	echo "$cache" | sed -n 1'p' | tr ',' '\n' | while read -r image; do
		echo "Pulling ${image}"
		docker pull "${image}" || true
	done
}

fetch_aperture_version() {
	component=aperture-controller
	tag_matcher="releases/${component}/*"
	current_branch="$(git branch --show-current)"

	if [[ "${current_branch}" == "stable/"* ]]; then
		version="$(git describe --match "${tag_matcher}" | cut -d/ -f 3)"
	else
		tag="$(git tag -l --sort="-version:refname" "${tag_matcher}" | head -n1 | cut -d/ -f 3)"
		commits="$(git rev-list "${tag}"..HEAD --count)"
		version="${tag##*/}-b.${commits}"
	fi

	echo "${version}" | cut -dv -f 2
}

if ! parse_tags_to_docker_arg; then
	echo "Unable to parse provided tags."
	echo "Check your \"tag\" parameter or refer to the docs and try again: https://circleci.com/developer/orbs/orb/circleci/docker."
	exit 1
fi

if [ -n "$PARAM_CACHE_FROM" ]; then
	if ! pull_images_from_cache; then
		echo "Unable to pull images from the cache."
		echo "Check your \"cache_from\" parameter or refer to the docs and try again: https://circleci.com/developer/orbs/orb/circleci/docker."
		exit 1
	fi
fi

build_args=(
	"--file=$PARAM_DOCKERFILE_PATH/$PARAM_DOCKERFILE_NAME"
	"$DOCKER_TAGS_ARG"
	"--build-arg=APERTURECTL_BUILD_GIT_BRANCH=$PARAM_GIT_BRANCH"
	"--build-arg=APERTURECTL_BUILD_GIT_COMMIT_HASH=$PARAM_GIT_COMMIT_HASH"
	"--build-arg=APERTURECTL_BUILD_VERSION=$(fetch_aperture_version)"
)

if [ -n "$PARAM_SSH_FORWARD" ]; then
	build_args+=("--ssh $PARAM_SSH_FORWARD")
fi

if [ -n "$PARAM_EXTRA_BUILD_ARGS" ]; then
	build_args+=("$PARAM_EXTRA_BUILD_ARGS")
fi

if [ -n "$PARAM_CACHE_FROM" ]; then
	build_args+=("--cache-from=$PARAM_CACHE_FROM")
fi

if [ "$PARAM_USE_BUILDKIT" -eq 1 ]; then
	build_args+=("--progress=plain")
fi

# The context must be the last argument.
build_args+=("$PARAM_DOCKER_CONTEXT")

old_ifs="$IFS"
IFS=' '

set -x
# expand env variables
argstring="$(eval echo "${build_args[*]}")"

# shellcheck disable=SC2048,SC2086
docker build $argstring
set +x

IFS="$old_ifs"
