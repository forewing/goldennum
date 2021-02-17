#!/bin/sh

set -x

rm -rf output
mkdir -p output

export GO111MODULE=on

# Build server
go build -ldflags "-s -w"
mv goldennum output/
