#!/bin/bash

3lv build \
    -s core \
    -f go.mod \
    -r ghcr.io/3lvia \
    --additional-tags latest,v0,alpha, \
    --go-main-package-directory . \
    cli
