#!/bin/bash

export GOPATH="$HOME/go:$PWD"
go build -o tracer ./src
