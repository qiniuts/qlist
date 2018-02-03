#!/usr/local/bin/bash

if [ -z "$GOPATH" ]
then
    export GOPATH=$(pwd)
else
    export GOPATH="${GOPATH}:"$(pwd)
fi
