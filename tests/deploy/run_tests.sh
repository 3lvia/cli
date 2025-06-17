#!/bin/bash

set -euo pipefail

test_disabled_deploy() {
    if 3lv deploy \
        -s core \
        -f tests/deploy/values.yaml \
        -i latest-cache \
        demo-api; then
        echo "Should fail since CI is not set"
        exit 1
    fi
}

# Must be signed in via az
test_aks_deploy() {
    if ! CI=true 3lv deploy \
        -s core \
        -f tests/deploy/values.yaml \
        -i latest-cache \
        --dry-run \
        demo-api; then
        echo "Failed to dry-run deploy to AKS"
        exit 1
    fi
}

main() {
    test_disabled_deploy
    # test_aks_deploy

    echo 'All tests passed!'
}

main
