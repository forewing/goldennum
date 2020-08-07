#!/bin/sh

rm -rf output
mkdir -p output

# Install packr2, go to / to prevent updating go.mod
(cd / && go get -u github.com/go-bindata/go-bindata/go-bindata)

export GO111MODULE=on

# Build resources files
go-bindata -fs -prefix "statics/" statics/ templates/

# Build server
go build -ldflags "-s -w"
mv goldennum output/
