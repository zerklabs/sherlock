#!/bin/bash
#
# This script builds the application from source.
# If we're building on Windows, specify an extension
EXTENSION=""
DISTPATH="bin/linux/"
if [ "$(go env GOOS)" = "windows" ]; then
    EXTENSION=".exe"
    DISTPATH="bin\\win\\"
fi


echo "--> Building sherlock-runner"
go build -v -o ${DISTPATH}sherlock-runner${EXTENSION} github.com/zerklabs/sherlock/runner
