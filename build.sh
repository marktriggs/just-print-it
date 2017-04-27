#!/bin/bash

set -e

cd "`dirname "$0"`"

export GOPATH=$PWD

TARGET=build/just_print_it

if [ "$1" = "dev" ]; then
    GOOS=linux GOARCH=amd64 go build -o $TARGET/just_print_it.linux just_print_it
    exit
fi

rm -rf build
mkdir -p $TARGET

GOOS=linux GOARCH=amd64 go build -o $TARGET/just_print_it.linux just_print_it
GOOS=darwin GOARCH=amd64 go build -o $TARGET/just_print_it.osx just_print_it

cp -a templates $TARGET

echo "build to $TARGET"
