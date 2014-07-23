#!/usr/bin/env bash

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR:$OLDGOPATH"
go build online
export GOPATH="$OLDGOPATH"
