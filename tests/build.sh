#!/bin/bash

3lv build \
    -s core \
    -f go.mod \
    -r ghcr.io/3lvia \
    --aditional-tags latest,v0,alpha, \
    --go-main-package-dir . \
    cli
