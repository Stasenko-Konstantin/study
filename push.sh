#! /bin/bash

rm иксис.zip
cd курсачи/иксис/archive
GOOS=windows GOARCH=amd64 go build main.go
cd ..
cd ..
cd ..
zip -r иксис.zip курсачи/иксис
