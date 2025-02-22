#!/bin/bash
rm -rf ./bin
env GOOS=linux GOARCH=amd64 go build -o ./bin/morpheus-fling
cd ./fling-decryptor/
rm -rf ./bin
env GOOS=darwin GOARCH=arm64 go build -o ./bin/fling-decryptor-osx
env GOOS=windows GOARCH=amd64 go build -o ./bin/fling-decryptor-windows
env GOOS=linux GOARCH=amd64 go build -o ./bin/fling-decryptor-linux
cd ..
openssl genrsa -traditional -out ./fling-decryptor/bin/redPill.pem 2048
openssl rsa -in ./fling-decryptor/bin/morpheus.pem -outform PEM -pubout -out ./bin/bluePill.pub
zip -r "morpheus-fling-$1.zip" ./bin/
zip -r "fling-decryptor-$1.zip" ./fling-decryptor/bin/
mv "morpheus-fling-$1.zip" ./releases
mv "fling-decryptor-$1.zip" ./releases

