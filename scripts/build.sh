#!/bin/bash

set -eou pipefail

main() {
    local out_dir="${1:-.}"
    local binary_path="$out_dir/3lv-linux-amd64"
    local compressed_package_path="$binary_path.tar.gz"

    export GO111MODULE=on
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64

    go build -o "$binary_path" ./cmd/3lv
    echo "$binary_path"

    if [[ "$2" == "--compress" ]]; then
        tar -czf \
            "$compressed_package_path" \
            "$binary_path" LICENSE VERSION README.md
    fi
}

main "$@"
