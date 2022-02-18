#!/bin/bash

cd "$(dirname $0)"
if [ -n "$(go fmt ./...)" ]; then
    echo "Go code is not formatted, run 'go fmt github.com/exercism/cli/...'" >&2
    exit 1
fi
