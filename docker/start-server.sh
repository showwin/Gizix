#!/bin/bash

# install packages
rm -rf /go/src/github.com/showwin/Gizix
mkdir -p /go/src/github.com/showwin
cp -r /go/src/development/showwin/Gizix /go/src/github.com/showwin/Gizix && echo 'Copy Packages: github.com/showwin'
go get -t -d -v ./... && echo 'Packages are up-to-date !'

# start Nginx
sudo service nginx start

# start MySQL
sudo service mysql start
mysql -u root mysql < docker/create_table.sql

# build and start
go build -o bin/application || exit
bin/application
