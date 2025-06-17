#!/bin/bash

set -euo pipefail

system_name='core'

test_create_dotnet8() {
    app_name='demo-api'
    project_dir='DemoApi' # pascal case for .NET

    for dotnet8_template_type in dotnet8-webapi dotnet8-worker; do
        output_dir="$(mktemp -d)"

        if ! 3lv create \
            -s "$system_name" \
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

        if [[ ! -f "$output_dir/$project_dir/.github/workflows/build-deploy-$app_name.yaml" ]]; then
            echo "Workflow file does not exist for template $dotnet8_template_type."
            exit 1
        fi

        current_dir="$(pwd)"
        cd "$output_dir/$project_dir"

        if ! 3lv build "$app_name"; then
            echo "Failed to build project for template $dotnet8_template_type."
            cd "$current_dir"
            exit 1
        fi

        cd "$current_dir"
    done
}

test_create_python() {
    app_name='demo-api-python'
    project_dir="$app_name"

    for python_template_type in python-webapi python-worker; do
        for python_version in 3.12 3.13; do
            output_dir="$(mktemp -d)"

            if ! 3lv create \
                -s "$system_name" \
                -a "$app_name" \
                -t "$python_template_type" \
                --python-version "$python_version" \
                --non-interactive \
                "$output_dir"; then
                echo "Failed to create project with template $python_template_type and version $python_version."
                exit 1
            fi

            if [[ ! -d "$output_dir/$project_dir" ]]; then
                echo "Project directory does not exist for template $python_template_type and version $python_version."
                exit 1
            fi

            if [[ $(cat "$output_dir/$project_dir/.python-version") != "$python_version" ]]; then
                echo "Python version file does not exist for template $python_template_type and version $python_version."
                exit 1
            fi

            if [[ ! -f "$output_dir/$project_dir/.github/workflows/build-deploy-$app_name.yaml" ]]; then
                echo "Workflow file does not exist for template $python_template_type and version $python_version."
                exit 1
            fi

            current_dir="$(pwd)"
            cd "$output_dir/$project_dir"

            if ! 3lv build "$app_name"; then
                echo "Failed to build project for template $python_template_type and version $python_version."
                cd "$current_dir"
                exit 1
            fi

            cd "$current_dir"
        done
    done
}

test_create_go() {
    app_name='demo-api-go'
    project_dir="$app_name"

    for go_template_type in go-webapi; do
        output_dir="$(mktemp -d)"

        if ! 3lv create \
            -s "$system_name" \
            -a "$app_name" \
            -t "$go_template_type" \
            --non-interactive \
            "$output_dir"; then
            echo "Failed to create project with template $go_template_type."
            exit 1
        fi

        if [[ ! -d "$output_dir/$project_dir" ]]; then
            echo "Project directory does not exist for template $go_template_type."
            exit 1
        fi

        if [[ ! -f "$output_dir/$project_dir/.github/workflows/build-deploy-$app_name.yaml" ]]; then
            echo "Workflow file does not exist for template $go_template_type."
            exit 1
        fi

        current_dir="$(pwd)"
        cd "$output_dir/$project_dir"

        if ! 3lv build "$app_name"; then
            echo "Failed to build project for template $go_template_type."
            cd "$current_dir"
            exit 1
        fi

        cd "$current_dir"
    done
}

main() {
    test_create_dotnet8
    test_create_python
    test_create_go

    echo 'All tests passed!'
}

main
