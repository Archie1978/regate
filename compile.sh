#!/bin/bash

go build

cd www/web-remotedektop
npm run build
cd ../..
./webRemotedektop
