#!/bin/bash

set -eou pipefail

main() {
    local binary_path="3lv-linux-amd64"
    local compressed_package_path="$binary_path.tar.gz"

    env GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$binary_path" ./cmd/3lv
    echo "$binary_path"

    if [[ "${1:-}" == '--compress' ]]; then
        tar -czf \
            "$compressed_package_path" \
            "$binary_path" LICENSE VERSION README.md
    fi
}

main "$@"
