#!/usr/bin/env bash

# this script checks github.com/fluxninja/aperture and returns the latest release version
# release versions follow semver, so we can use the github api to get the latest release
# ignoring pre-releases

# get the latest release version
latest_release=$(curl -s https://api.github.com/repos/fluxninja/aperture/releases/latest | grep tag_name | cut -d '"' -f 4)

# remove the leading 'v' from the version
latest_release=${latest_release:1}

# print the version
echo "$latest_release"
