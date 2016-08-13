#!/bin/bash

CONTAINER=go-dotnet
WORKDIR=/go/src/github.com/matiasinsaurralde/go-dotnet

docker build -t $CONTAINER .

docker run -it --rm -v `pwd`/../:$WORKDIR -w $WORKDIR  $CONTAINER 
