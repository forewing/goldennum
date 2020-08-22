#!/bin/sh

set -x

rm -rf output
mkdir -p output

# Install packr2, go to / to prevent updating go.mod
(cd / && go get -u github.com/go-bindata/go-bindata/go-bindata)

export GO111MODULE=on

# Build resources files
go-bindata -fs -prefix "statics/" statics/ templates/

# Fix go-bindata's BUG on windows with prefix\
case "$(uname -s)" in
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        sed -i 's/templates\\[^\\]/templates\//g' bindata.go
        sed -i 's/templates\\\\/templates\//g' bindata.go
        ;;
    *)
        ;;
esac

# Build server
go build -ldflags "-s -w"
mv goldennum output/
