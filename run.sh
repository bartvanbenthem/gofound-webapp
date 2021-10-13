#!/bin/bash
go build -o build/bin/webserver $(ls -1 cmd/web/*.go | grep -v _test.go)
./build/bin/webserver