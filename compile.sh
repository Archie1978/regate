#!/bin/bash

#
#
#  Example compilation linux and windows
#
#

cd webservice/www/regate/
npm run build

if [ $? -ne 0 ]; then
echo "Package not found, install all package"
npm install
npm install xterm-addon-web-links
npm run build
if [ $? -ne 0 ]; then
exit 1
fi
fi

cd ../../..


cd cmd/regate-standalone-user/

if [ $1 -eq "win" ]; then
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build
else
go build 
fi
