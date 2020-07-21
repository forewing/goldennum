#!/bin/bash

rm -rf output
mkdir -p output

go build
mv goldennum output/

cp -r templates output/
cp -r statics output/
