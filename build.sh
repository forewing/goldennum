#!/bin/sh

rm -rf output
mkdir -p output

# Install packr2, go to / to prevent updating go.mod
(cd / && go get github.com/gobuffalo/packr/v2/packr2)

export GO111MODULE=on

# Build resources files
packr2

# Build server
go build -ldflags "-s -w"
mv goldennum output/
