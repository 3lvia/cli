#!/bin/bash

set -e

test_normal_scan() {
    if 3lv scan debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_table_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats table \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_json_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats json \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "File 'trivy.json' not found"
        exit 1
    fi
}

test_sarif_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats sarif \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_markdown_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats markdown \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.md ]]; then
        echo "ERROR: File 'trivy.md' not found"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_json_sarif_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats json,sarif \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not found"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi
}

test_json_markdown_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats json,markdown \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not found"
        exit 1
    fi

    if [[ ! -f trivy.md ]]; then
        echo "ERROR: File 'trivy.md' not found"
        exit 1
    fi
}

test_json_table_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats json,table \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not found"
        exit 1
    fi
}

test_sarif_markdown_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats sarif,markdown \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi

    if [[ ! -f trivy.md ]]; then
        echo "ERROR: File 'trivy.md' not found"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_sarif_table_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats sarif,table \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_json_sarif_table_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats json,sarif,table \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not found"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi
}

test_sarif_table_markdown_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats sarif,table,markdown \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi

    if [[ ! -f trivy.md ]]; then
        echo "ERROR: File 'trivy.md' not found"
        exit 1
    fi

    if [[ -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not cleaned up"
        exit 1
    fi
}

test_all_outputs_scan() {
    if 3lv scan \
        --severity CRITICAL,HIGH \
        --formats table,json,sarif,markdown \
        debian:10; then
        echo "Scan should have failed"
        exit 1
    fi

    if [[ ! -f trivy.json ]]; then
        echo "ERROR: File 'trivy.json' not found"
        exit 1
    fi

    if [[ ! -f trivy.sarif ]]; then
        echo "ERROR: File 'trivy.sarif' not found"
        exit 1
    fi

    if [[ ! -f trivy.md ]]; then
        echo "ERROR: File 'trivy.md' not found"
        exit 1
    fi
}

cleanup_files() {
    if [[ -f trivy.json ]]; then
        echo 'Removing trivy.json'
        rm trivy.json
    fi

    if [[ -f trivy.sarif ]]; then
        echo 'Removing trivy.sarif'
        rm trivy.sarif
    fi

    if [[ -f trivy.md ]]; then
        echo 'Removing trivy.md'
        rm trivy.md
    fi
}


main() {
    test_normal_scan
    test_table_scan
    test_json_scan
    test_sarif_scan
    test_markdown_scan
    test_json_sarif_scan
    test_json_markdown_scan
    test_json_table_scan
    test_sarif_markdown_scan
    test_sarif_table_scan
    test_json_sarif_table_scan
    test_sarif_table_markdown_scan
    test_all_outputs_scan

    cleanup_files

    echo 'All tests passed!'
}

main
