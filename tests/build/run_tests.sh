#!/bin/bash

set -e

test_build_cli() {
    if ! 3lv build \
        -s core \
        -f go.mod \
        -r ghcr.io/3lvia \
        --additional-tags latest,v0,alpha, \
        --go-main-package-directory . \
        cli; then
        echo "Failed to build CLI"
        exit 1
    fi
}

test_build_dockerfile() {
    if 3lv build \
        -s core \
        -f tests/build/Dockerfile \
        -r ghcr.io/3lvia \
        vulnerable-service; then
        echo "Did not exit with error, should fail due to vulnerabilities in base image"
        exit 1
    fi
}

main() {
    test_build_cli
    test_build_dockerfile

    echo 'All tests passed!'
}

main
