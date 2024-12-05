#!/bin/bash

set -e

test_create_dotnet8() {
    app_name='demo-api'
    project_dir='DemoApi' # pascal case for .NET

    for dotnet8_template_type in dotnet8-webapi dotnet8-worker; do
        output_dir="$(mktemp -d)"

        if ! 3lv create \
            -s core \
            -a "$app_name" \
            -t "$dotnet8_template_type" \
            --non-interactive \
            "$output_dir"; then
            echo "Failed to create project with template $dotnet8_template_type."
            exit 1
        fi

        if [[ ! -d "$output_dir/$project_dir" ]]; then
            echo "Project directory does not exist for template $dotnet8_template_type."
            exit 1
        fi

        if [[ ! -f "$output_dir/$project_dir/.github/workflows/build-deploy-$app_name.yml" ]]; then
            echo "Workflow file does not exist for template $dotnet8_template_type."
            exit 1
        fi

        # TODO: use 3lv build
        if ! dotnet build "$output_dir/$project_dir"; then
            echo "Failed to build project for template $dotnet8_template_type."
            exit 1
        fi
    done
}

main() {
    test_create_dotnet8

    echo 'All tests passed!'
}

main
