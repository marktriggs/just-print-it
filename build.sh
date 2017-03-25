#!/bin/bash

set -e

cd "`dirname "$0"`"

export GOPATH=$PWD

rm -rf build

TARGET=build/just_print_it

mkdir -p $TARGET

GOOS=linux GOARCH=386 go build -o $TARGET/just_print_it.linux just_print_it
GOOS=darwin GOARCH=386 go build -o $TARGET/just_print_it.osx just_print_it

cp -a templates $TARGET

echo "build to $TARGET"
