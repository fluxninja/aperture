#!/bin/bash

# Based on https://github.com/giantswarm/grumpy

set -o errexit
set -o nounset
set -o pipefail

# CREATE THE PRIVATE KEY FOR OUR CUSTOM CA
openssl genrsa -out ca.key 2048

# GENERATE A CA CERT WITH THE PRIVATE KEY
openssl req -new -x509 -key ca.key -out ca.crt -config ca_config.txt -days 3650

# CREATE THE PRIVATE KEY FOR OUR GRUMPY SERVER
openssl genrsa -out key.pem 2048

# CREATE A CSR FROM THE CONFIGURATION FILE AND OUR PRIVATE KEY
openssl req -new -key key.pem -subj "/CN=agent-webhooks.aperture-system.svc" -out csr.csr \
    -config config.txt -days 3650

# CREATE THE CERT SIGNING THE CSR WITH THE CA CREATED BEFORE
openssl x509 -req -in csr.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out crt.pem \
    -extensions v3_req -extfile config.txt -days 3650

# # INJECT CA IN THE WEBHOOK CONFIGURATION
sed -i -E -e 's/(caBundle: ").*(" # agent-cm-validator)$/\1'"$(base64 -w 0 ca.crt)"'\2/' \
    "$(git rev-parse --show-toplevel)/ops/charts/agent/templates/webhooks.yaml"
