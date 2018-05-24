#!/bin/bash
export GOPATH=/c/gopath
cd $GOPATH/src/github.com/matiasinsaurralde/go-dotnet/dotnet
go get -u
go test -v