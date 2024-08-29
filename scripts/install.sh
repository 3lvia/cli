#!/bin/bash

main() {
    local binary_path
    if [[ "$1" == "--build" ]]; then
        binary_path=$(./build.sh)
    else
        binary_path="$1"
    fi

    local binary_name
    IFS=- read binary_name _ <<< "$binary_path"
    mv "$binary_path" "$binary_name"

    sudo install -Dm755 -t /usr/bin "$binary_name"
}

main "$@"
