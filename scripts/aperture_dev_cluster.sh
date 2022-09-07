#!/usr/bin/env bash

set -euo pipefail

# Values selected due to minikube using them by default
# If this is changed, update entries in GCP
readonly cluster_docker_network="kind"
git_root="$(git rev-parse --show-toplevel)"
readonly git_root="${git_root}"
readonly kind_node_version='kindest/node:v1.23.3'
readonly registry_name='docker-registry'
readonly registry_port='5001' # Docker registry usually uses port 5000, but MacOS Monterey uses it for AirPlay server
readonly registry_image='registry:2'

setup() {
  create_registry
  case "${1:-}" in
    kind) setup_kind_cluster;;
    *) printf 'Unknown cluster provider - "%s"\n' "${1:-}"; return 1;;
  esac
  document_registry
}

teardown() {
  case "${1:-}" in
    kind) teardown_kind_cluster;;
    *) printf 'Unknown cluster provider - "%s"\n' "${1:-}"; return 1;;
  esac
}

main() {
  local -r subcommand="${1:-}"
  local -r cluster_provider="${2:-${NINJA_DEV_CLUSTER_PROVIDER:-kind}}"
  # We need to handle it earlier, before cluster_provider check logic kicks in
  if [ "${subcommand:-}" = "help" ]; then
    print_help help
    return
  fi

  # Unused for now, but we might use it soon-ish if we add rancher support for MacOS folks
  if [ -z "${cluster_provider:-}" ]; then
    printf 'ERROR: You need to specify what cluster provider to use.\n' >&2
    printf 'See "help" subcommand for details.\n' >&2
    return 1
  fi

  if ! command -v "${cluster_provider}" &>/dev/null; then
    # The provider command is not available
    printf 'ERROR: The selected provider is not installed. Install it by going to:\n' >&2
    case "${cluster_provider}" in
      kind) print 'https://kind.sigs.k8s.io/docs/user/quick-start/#installation\n' >&2;;
      *) printf 'Unknown cluster provider - "%s"\n' "${cluster_provider}";;
    esac
    return 1
  fi

  case "${subcommand:-}" in
    up|start|setup|create) setup "${cluster_provider}";;
    down|stop|teardown|delete) teardown "${cluster_provider}";;
    flip|redo|recreate) teardown "${cluster_provider}" && setup "${cluster_provider}";;
    *)
      print_help "${subcommand:-}"
      return 1
      ;;
  esac
}


#########################
### CLUSTER FUNCTIONS ###
#########################

setup_kind_cluster() {
  start_kind
  connect_registry_to_kind_network
}

teardown_kind_cluster() {
  disconnect_registry_from_kind_network
  kind delete cluster
}

########################
### HELPER FUNCTIONS ###
########################

print_help() {
  local -r passed_argument="${1:-}"
  local -r argument_msg_suffix='Should be either "up", "down" or "flip"'
  if [ "${passed_argument:-}" = "help" ]; then
    : # Don't print any warning/error
  elif [ -z "${passed_argument:-}" ]; then
    printf 'Missing argument. %s\n' "${argument_msg_suffix}" >&2
  else
    printf 'Invalid argument "%s". %s\n' "${passed_argument}" "${argument_msg_suffix}" >&2
  fi
  local -r script_called_name="${BASH_SOURCE[0]}"
  printf 'Usage: %s <up|down|flip> [kind]\n' "${script_called_name}" >&2
  printf '\tYou can set NINJA_DEV_CLUSTER_PROVIDER env var to declare default value for second argument.\n' >&2
  printf '\tIf neither argument is passed or envar set, it defaults to "kind".\n' >&2
}

create_registry() {
  # create registry container unless it already exists
  if [ "$(docker inspect -f '{{.State.Running}}' "${registry_name}" 2>/dev/null || true)" != 'true' ]; then
    printf 'Creating new docker registry. It will not be deleted automatically at any point.\n'
    printf 'To manually remove it, run:\n\tdocker rm --force %s\n' "${registry_name}"
    local docker_run_args=(
      --detach
      --restart"="always
      --publish "127.0.0.1:${registry_port}:5000"
      --name "${registry_name}"
      "${registry_image}"
    )
    docker run "${docker_run_args[@]}"
  fi
}

is_connected_to_network() {
  local -r container_name="${1?}"
  local -r network_name="${2?}"
  local networks_str
  networks_str="$(docker inspect -f='{{json .NetworkSettings.Networks}}' "${container_name}" | jq -r 'keys|.[]')"
  readarray -t networks <<<"${networks_str}"
  for network in "${networks[@]}"; do
    if [ "${network}" = "${network_name}" ]; then
      return 0
    fi
  done
  return 1
}

connect_registry_to_kind_network() {
  # connect the registry to the cluster network if not already connected
  if ! is_connected_to_network "${registry_name}" "${cluster_docker_network}"; then
    printf 'Connecting registry to kind network...\n'
    docker network connect "${cluster_docker_network}" "${registry_name}"
  fi
}

disconnect_registry_from_kind_network() {
  if is_connected_to_network "${registry_name}" "${cluster_docker_network}"; then
    docker network disconnect "${cluster_docker_network}" "${registry_name}"
  fi
}

document_registry() {
  # Document the local registry
  # https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
  cat <<EOF | kubectl apply -f -
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: local-registry-hosting
    namespace: kube-public
  data:
    localRegistryHosting.v1: |
      host: "localhost:${registry_port}"
      help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
}

start_kind() {
  if kind get clusters | grep --quiet --fixed-strings --line-regexp kind >/dev/null; then
    printf 'WARN: kind is already running - not recreating\n' >&2
  else
    KIND_EXPERIMENTAL_DOCKER_NETWORK="${cluster_docker_network}" kind create cluster --image "${kind_node_version}" --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
featureGates:
  ServiceInternalTrafficPolicy: true
nodes:
- role: control-plane
- role: worker
- role: worker
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${registry_port}"]
    endpoint = ["http://${registry_name}:5000"]
EOF
  fi
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
