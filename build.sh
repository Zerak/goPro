#!/usr/bin/env bash

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

rm -r $GOPATH/pkg
rm -r $GOPATH/bin
mkdir $GOPATH/bin
cd $GOPATH/bin

echo $GOPATH

go fmt ../src/

go build ../src/

export GOPATH="$OLDGOPATH"

echo 'build finished'
