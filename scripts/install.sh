#!/bin/bash

main() {
    local binary_path
    binary_path=$(./build.sh)
    local binary_name
    IFS=- read binary_name _ <<< "$binary_path"
    mv "$binary_path" "$binary_name"

    sudo install -Dm755 -t /usr/bin "$binary_name"
}

main
