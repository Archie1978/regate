#!/bin/bash


cd webservice/www/regate/
npm run build
cd ../../..


cd cmd/regate-standalone-user/
go build 
