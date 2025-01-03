#!/bin/bash

set -euo pipefail

# TODO: find a way to test the init command without user interaction

test_init_non_interactive() {
    if 3lv init --non-interactive; then
        echo 'Should fail when running init non-interactively.'
        exit 1
    fi
}

main() {
    test_init_non_interactive

    echo 'All tests passed!'
}

main
