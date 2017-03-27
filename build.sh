#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN

go build src/mysite/main.go

if [ ! -d ./bin ]; then
	mkdir bin
fi

if [ -e ./main ]; then
   mv main ./bin/
fi

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'build finished'

