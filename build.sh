#!/bin/sh

rm -rf output
mkdir -p output

go build -ldflags "-s -w"
mv goldennum output/

cp -r templates output/
cp -r statics output/
