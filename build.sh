#!/bin/bash
set -e
set -x

if [ ! -d /tmp/build-golang/src/github.com/szyhf ]; then
    mkdir -p /tmp/build-golang/src/github.com/szyhf
    ln -s $PWD /tmp/build-golang/src/github.com/szyhf/go-jsoniter
fi
export GOPATH=/tmp/build-golang
go get -u github.com/golang/dep/cmd/dep
cd /tmp/build-golang/src/github.com/szyhf/go-jsoniter
exec $GOPATH/bin/dep ensure -update
