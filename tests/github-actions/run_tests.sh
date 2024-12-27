#!/bin/bash

set -e

test_github_actions() {
    app_name='demo-api'
    output_dir=$(mktemp -d)

    if ! 3lv gha \
        -s core \
        -a "$app_name" \
        -f "$app_name.csproj" \
        "$output_dir"; then
        echo "Failed to add GitHub Actions."
        exit 1
    fi

    if [[ ! -f "$output_dir/.github/workflows/build-deploy-$app_name.yml" ]]; then
        echo 'Workflow file does not exist.'
        exit 1
    fi
}

main() {
    test_github_actions

    echo 'All tests passed!'
}

main
