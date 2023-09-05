#!/bin/bash

go build

cd www/regate
npm run build
cd ../..
./regate
