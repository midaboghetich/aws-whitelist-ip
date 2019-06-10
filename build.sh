#!/bin/bash

if [ -z "$1" ]
then
      echo "you need to define a version number to build this application"
      exit 2
fi

version=$1
time=$(date)

rm -rf linux
rm -rf osx
rm -rf win

env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" .
mkdir osx 
mv aws-whitelist-ip ./osx/

env GOOS=windows GOARCH=amd64 go  build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" .
mkdir win
mv aws-whitelist-ip.exe ./win/

env GOOS=linux GOARCH=amd64 go  build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" .
mkdir linux
mv aws-whitelist-ip ./linux/