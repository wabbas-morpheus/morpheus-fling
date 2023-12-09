#!/bin/bash
env GOOS=darwin GOARCH=arm64 go build -o ./bin/morpheus-fling-osx
env GOOS=windows GOARCH=amd64 go build -o ./bin/morpheus-fling-windows
env GOOS=linux GOARCH=amd64 go build -o ./bin/morpheus-fling-linux