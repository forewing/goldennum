#!/bin/sh

rm -rf output
mkdir -p output

VERSION_PACKAGE="github.com/forewing/goldennum/version"
LDFLAGS="-s -w"

GIT_TAG=$(git describe --tags)
if [ $? -eq 0 ]; then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Version=$GIT_TAG'"
fi

GIT_HASH=$(git rev-parse HEAD)
if [ $? -eq 0 ]; then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Hash=$GIT_HASH'"
fi

# Build server
go build -ldflags "$LDFLAGS"
mv goldennum output/
