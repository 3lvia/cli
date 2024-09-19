#!/bin/bash

set -e

# Must be signed in via az
test_aks_deploy() {
    if ! 3lv deploy \
        -s core \
        -f tests/deploy/values.yml \
        -i latest-cache \
        --skip-authentication \
        --dry-run \
        demo-api; then
        echo "Failed to dry-run deploy to AKS"
        exit 1
    fi
}

main() {
    test_aks_deploy

    echo 'All tests passed!'
}

main
