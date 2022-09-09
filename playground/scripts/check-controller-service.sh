#!/usr/bin/env bash
# Port-forward using kubectl, with simple healthchecks to work around bug in kubectl
set -euo pipefail

while (kubectl get svc aperture-controller -n aperture-controller); ret=$?; [ $ret -eq 0 ]; do
  echo service available; sleep 2;
done
