#!/usr/bin/env bash
set -ex
export GOBIN=$GOPATH/bin

go get
go get github.com/konsorten/go-windows-terminal-sequences
go test
make
