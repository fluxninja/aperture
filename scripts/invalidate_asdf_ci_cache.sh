#!/usr/bin/env bash

FIND="find"
# use gawk on linux and mac
AWK="gawk"
if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

# for all yml or yaml files in .circleci directory, increment asdf cache version which looks like aperture-asdf-cache-v<version>-*. it should look like aperture-asdf-cache-v<version+1>-*
# shellcheck disable=SC2016
"$FIND" .circleci -type f \( -name "*.yml" -o -name "*.yaml" \) -print0 | while read -r -d $'\0' file; do
	"$AWK" -i inplace '/aperture-asdf-cache-v[0-9]+-/ {match($0, /aperture-asdf-cache-v([0-9]+)(.*)/, a); $0 = gensub(/aperture-asdf-cache-v[0-9]+-/, "aperture-asdf-cache-v" a[1]+1 "-", "g", $0)} {print}' "$file"
done

# invalidate aperture-v<version>-daily-cache-* as well
# shellcheck disable=SC2016
"$FIND" .circleci -type f \( -name "*.yml" -o -name "*.yaml" \) -print0 | while read -r -d $'\0' file; do
	"$AWK" -i inplace '/aperture-v[0-9]+-daily-cache-/ {match($0, /aperture-v([0-9]+)(.*)/, a); $0 = gensub(/aperture-v[0-9]+-daily-cache-/, "aperture-v" a[1]+1 "-daily-cache-", "g", $0)} {print}' "$file"
done

# invalidate aperture-v<version>-go-cache-* as well
# shellcheck disable=SC2016
"$FIND" .circleci -type f \( -name "*.yml" -o -name "*.yaml" \) -print0 | while read -r -d $'\0' file; do
	"$AWK" -i inplace '/aperture-v[0-9]+-go-cache-/ {match($0, /aperture-v([0-9]+)(.*)/, a); $0 = gensub(/aperture-v[0-9]+-go-cache-/, "aperture-v" a[1]+1 "-go-cache-", "g", $0)} {print}' "$file"
done
