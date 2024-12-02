#!/bin/bash

set -e

test_create_dotnet() {
    output_dir="$(mktemp -d)"
    app_name='demo-api'
    project_dir='DemoApi' # pascal case for .NET

    if ! 3lv create \
        -s core \
        -a "$app_name" \
        -t dotnet \
        "$output_dir"; then
        echo "Failed to create project"
        exit 1
    fi

    if [[ ! -d "$output_dir/$project_dir" ]]; then
        echo "Project directory does not exist"
        exit 1
    fi

    if [[ ! -f "$output_dir/$project_dir/.github/workflows/build-deploy-$app_name.yml" ]]; then
        echo "Workflow file does not exist"
        exit 1
    fi
}

main() {
    test_create_dotnet

    echo 'All tests passed!'
}

main
