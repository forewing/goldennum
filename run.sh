#!/bin/bash

export $(grep -v '^#' .env | xargs)

./output/goldennum
