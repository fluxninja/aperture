#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

REPO_URL="https://x-access-token:${GITHUB_TOKEN}@github.com/fluxninja/aperture.git"
BRANCH="gh-pages"
TARGET_DIR="."
INDEX_DIR="."
CHARTS_URL="https://fluxninja.github.io/aperture/"
COMMIT_EMAIL="${CIRCLE_PROJECT_USERNAME}@users.noreply.github.com"

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
      CHARTS+=("${dir}")
      echo "Found chart directory ${dir}"
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

package() {
  for chart in "${CHARTS[@]}"; do
    CHART_VERSION_CMD=$(helm inspect chart "$chart" | grep version | cut -d' ' -f2|tr -d " \t\n\r")
    APP_VERSION_CMD=$(helm inspect chart "$chart" | grep appVersion | cut -d' ' -f2|tr -d " \t\n\r")
    helm package "$chart" --destination "${CHARTS_TMP_DIR}" --version "$CHART_VERSION_CMD" --app-version "$APP_VERSION_CMD"
  done
}

upload(){
  tmpDir=$(mktemp -d)
  pushd "$tmpDir" >& /dev/null

  git clone "${REPO_URL}"
  cd aperture
  git config user.name "${CIRCLE_PROJECT_USERNAME}"
  git config user.email "${COMMIT_EMAIL}"
  git remote set-url origin "${REPO_URL}"
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
