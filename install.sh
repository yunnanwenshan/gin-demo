#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN

#go get github.com/garyburd/redigo/redis
#go get github.com/gorilla/context
#go get github.com/gorilla/securecookie
go get github.com/gorilla/sessions

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'build finished'

