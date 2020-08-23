#!/bin/sh

set -x

rm -rf output
mkdir -p output

export GO111MODULE=on

# Generate resources files
go generate -x

# Build server
go build -ldflags "-s -w"
mv goldennum output/
