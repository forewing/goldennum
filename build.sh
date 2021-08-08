#!/bin/sh

set -e

rm -rf output
mkdir -p output

OUTPUT="goldennum"

VERSION_PACKAGE="github.com/forewing/goldennum/version"
LDFLAGS="-s -w"

if GIT_TAG=$(git describe --tags); then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Version=$GIT_TAG'"
    OUTPUT="${OUTPUT}-${GIT_TAG}"
fi

if GIT_HASH=$(git rev-parse HEAD); then
    LDFLAGS="$LDFLAGS -X '$VERSION_PACKAGE.Hash=$GIT_HASH'"
fi

CMD_BASE="CGO_ENABLED=1 go build -trimpath -ldflags \"${LDFLAGS}\""

if [ ! -n "$1" ] || [ ! $1 = "all" ]; then
    eval ${CMD_BASE}
    mv goldennum output/
    exit 0
fi

# Cross compile

if [ ! $(uname) = "Linux" ] || [ -z $(which x86_64-w64-mingw32-gcc) ]; then
    echo You need Linux environment with x86_64-w64-mingw32-gcc installed to cross compile
    exit 1
fi

compress_tar_gz(){
    tar caf "${1}.tar.gz" "${1}"
    mv "${1}.tar.gz" output/
    rm "${1}"
}

compress_zip(){
    zip -q -r "${1}.zip" "${1}.exe"
    mv "${1}.zip" output/
    rm "${1}.exe"
}

# Linux
OUTPUT_FULL="${OUTPUT}-linux-amd64"
CMD="GOOS=linux GOARCH=amd64 ${CMD_BASE} -o ${OUTPUT_FULL}"
eval $CMD
compress_tar_gz "${OUTPUT_FULL}"

# Windows
OUTPUT_FULL="${OUTPUT}-windows-amd64"
CMD="CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 ${CMD_BASE} -o ${OUTPUT_FULL}.exe"
eval $CMD
compress_zip "${OUTPUT_FULL}"
