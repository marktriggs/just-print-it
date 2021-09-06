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

GOOS=linux GOARCH=amd64 go build -o $TARGET/just_print_it.linux src/just_print_it/main.go
GOOS=darwin GOARCH=arm64 go build -o $TARGET/just_print_it.osx src/just_print_it/main.go

cp -a templates $TARGET

echo "build to $TARGET"
