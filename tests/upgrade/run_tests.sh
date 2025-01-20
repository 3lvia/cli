#!/bin/bash

set -eou pipefail

test_upgrade() {
    if ! 3lv upgrade --force-reinstall; then
        echo 'Failed to upgrade 3lv.'
        exit 1
    fi
}

main() {
    test_upgrade

    echo 'All tests passed!'
}

main

