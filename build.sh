#!/bin/sh

rm -rf output
mkdir -p output

VERSION_PACKAGE="github.com/forewing/goldennum/version"
LDFLAGS="-s -w"

if GIT_TAG=$(git describe --tags); then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Version=$GIT_TAG'"
fi

if GIT_HASH=$(git rev-parse HEAD); then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Hash=$GIT_HASH'"
fi

# Build server
go build -ldflags "$LDFLAGS"
mv goldennum output/
