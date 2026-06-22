#!/usr/bin/env bash
set -e
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

make kubevirtci_init
make cluster_up
make cluster_sync_kubevirt
export KUBECONFIG=$(${CURRENT_DIR}/../kubeconfig.sh)
export IMAGE_REGISTRY="registry:5000"
export PUSH_REGISTRY="localhost:$(${CURRENT_DIR}/../cluster/cli.sh ports registry | tr -d '\r')"
export REQUIRE_IMAGE_PUSH_TLS_VERIFICATION=false
make image_test_dependencies
make push_test_dependencies
make sync_test_dependencies
make images
make push
make sync
