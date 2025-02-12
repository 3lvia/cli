#!/bin/bash

set -euo pipefail

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

    if ! [[ -f /tmp/3lv-cli-output/image-name ]]; then
        echo "Output image-name is not set"
        exit 1
    fi
}

test_build_dockerfile() {
    if 3lv build \
        -s core \
        -f tests/build/Dockerfile.vulnerable \
        -r ghcr.io/3lvia \
        vulnerable-service; then
        echo "Did not exit with error, should fail due to vulnerabilities in base image"
        exit 1
    fi
}

test_disable_scan_error() {
    if ! 3lv build \
        -s core \
        -f tests/build/Dockerfile.vulnerable \
        -r ghcr.io/3lvia \
        --scan-disable-error \
        vulnerable-service; then
        echo "Should not fail due to vulnerabilities in base image"
        exit 1
    fi

    if ! [[ -f /tmp/3lv-cli-output/image-name ]]; then
        echo "Output image-name is not set"
        exit 1
    fi
}

test_build_with_build_args() {
    # We disable scan errors in case the base image has vulnerabilities,
    # which we don't care about for this test.

    local environment='dev'
    local version='0.1.0'

    if ! 3lv build \
        -s core \
        -f tests/build/Dockerfile.build-args \
        --build-args "ENVIRONMENT=$environment,VERSION=$version" \
        --scan-disable-error \
        --go-main-package-directory . \
        cli; then
        echo "Failed to build CLI"
        exit 1
    fi

    if ! [[ -f /tmp/3lv-cli-output/image-name ]]; then
        echo "Output image-name is not set"
        exit 1
    fi

    local image_name
    image_name=$(cat /tmp/3lv-cli-output/image-name)

    local docker_run_output
    docker_run_output=$(docker run --rm "$image_name")

    local docker_run_output_expected="$environment $version"

    if [[ "$docker_run_output" != "$docker_run_output_expected" ]]; then
        echo "Build args not set correctly: got $docker_run_output, expected $docker_run_output_expected"
        exit 1
    fi
}

main() {
    test_build_cli
    test_build_dockerfile
    test_disable_scan_error
    test_build_with_build_args

    echo 'All tests passed!'
}

main
