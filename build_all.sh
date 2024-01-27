#!/bin/bash
env GOOS=darwin GOARCH=arm64 go build -o ./bin/morpheus-fling-osx
env GOOS=windows GOARCH=amd64 go build -o ./bin/morpheus-fling-windows
env GOOS=linux GOARCH=amd64 go build -o ./bin/morpheus-fling-linux
cd ./fling-decryptor/
env GOOS=darwin GOARCH=arm64 go build -o ./bin/fling-decryptor-osx
env GOOS=windows GOARCH=amd64 go build -o ./bin/fling-decryptor-windows
env GOOS=linux GOARCH=amd64 go build -o ./bin/fling-decryptor-linux
cd ..
openssl genrsa -traditional -out ./fling-decryptor/bin/morpheus.pem 2048
openssl rsa -in ./fling-decryptor/bin/morpheus.pem -outform PEM -pubout -out ./bin/morpheus.pub
