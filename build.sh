#!/bin/bash

if [ -z "$GOPATH" ]; then
    echo "Need to set GOPATH see http://golang.org for more info"
    exit 1
fi

# Install third party packages
go get github.com/gorilla/mux

# Compile and install
go install github.com/sjhitchner/sourcegraph 
