#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

REPO_URL="git@github.com:fluxninja/aperture.git"
BRANCH="gh-pages"
TARGET_DIR="."
INDEX_DIR="."
CHARTS_URL="https://fluxninja.github.io/aperture/"
COMMIT_EMAIL="ops@fluxninja.com"

CHARTS=()
CHARTS_TMP_DIR=$(mktemp -d)

main(){
  install_helm
  locate
  dependencies
  package
  upload
}

install_helm(){
  curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
}

locate() {
  dirs=$(find manifests/charts -name 'values.yaml' -exec dirname {} \;)
  for dir in ${dirs}; do
    if [[ -f "${dir}/Chart.yaml" ]]; then
      tag_name=$(echo "${CIRCLE_TAG}" | cut -d "/" -f 3)
      dir_name=$(echo "${dir}" | cut -d "/" -f 3)
      if [[ $tag_name == "${dir_name}" ]]; then
      CHARTS+=("${dir}")
      echo "Found chart directory ${dir}"
      fi
    else
      echo "Ignoring non-chart directory ${dir}"
    fi
  done
}


dependencies() {
  for chart in "${CHARTS[@]}"; do
    helm dependency update "${chart}"
  done
}

version_check() {
# Refer the issue for the Regex https://github.com/semver/semver/issues/232
rx='^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$'
if [[ $VERSION_CMD =~ $rx ]]; then
 echo "INFO:<-->Version $VERSION_CMD"
else
 echo "ERROR:<->Unable to validate package version: '$VERSION_CMD'"
 exit 1
fi
}

package() {
  for chart in "${CHARTS[@]}"; do
    # For now we are using the same version for app-version and version. Change it to two if we are using different versions.
    VERSION_CMD=$(echo "$CIRCLE_TAG" | cut -d "/" -f 4)
    if version_check "$VERSION_CMD"
    then
        helm package "$chart" --destination "${CHARTS_TMP_DIR}" --version "$VERSION_CMD" --app-version "$VERSION_CMD"
    fi
  done
}

upload(){
  tmpDir=$(mktemp -d)
  pushd "$tmpDir" >& /dev/null

  git clone "${REPO_URL}"
  cd aperture
  git config user.name "${CIRCLE_PROJECT_USERNAME}"
  git config user.email "${COMMIT_EMAIL}"
  git checkout "${BRANCH}"

  charts=$(cd "${CHARTS_TMP_DIR}" && find . -print0 | xargs -0)

  mkdir -p ${TARGET_DIR}

  if [[ -f "${INDEX_DIR}/index.yaml" ]]; then
    echo "Found index, merging changes"
    helm repo index "${CHARTS_TMP_DIR}" --url ${CHARTS_URL} --merge "${INDEX_DIR}/index.yaml"
    mv -f "${CHARTS_TMP_DIR}"/*.tgz ${TARGET_DIR}
    mv -f "${CHARTS_TMP_DIR}"/index.yaml ${INDEX_DIR}/index.yaml
  else
    echo "No index found, generating a new one"
    mv -f "${CHARTS_TMP_DIR}"/*.tgz ${TARGET_DIR}
    helm repo index ${INDEX_DIR} --url ${CHARTS_URL}
  fi

  git add ${TARGET_DIR}
  git add ${INDEX_DIR}/index.yaml

  git commit -m "Publish $charts"
  git push origin ${BRANCH}

  popd >& /dev/null
  rm -rf "$tmpDir"
}

main
