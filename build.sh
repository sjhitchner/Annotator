#!/bin/bash

if [ -z "$GOPATH" ]; then
    echo "Need to set GOPATH see http://golang.org for more info"
    exit 1
fi

# Install third party packages
echo "Installing Dependencies"
go get github.com/gorilla/mux

# Compile and install
echo "Building source..."
go install github.com/sjhitchner/annotator

echo "Run $GOPATH/bin/annotator"
