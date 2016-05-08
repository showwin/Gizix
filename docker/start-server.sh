#!/bin/bash

# install packages
rm -rf /go/src/github.com/showwin/Gizix
cp -r /go/src/development/showwin/Gizix /go/src/github.com/showwin/Gizix && echo 'Copy Packages: github.com/showwin'
go get -t -d -v ./... && echo 'Packages are up-to-date !'

# build and start
go build -o bin/application || exit
bin/application
