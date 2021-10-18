#!/bin/bash

# go test -v
# go test -cover
# go test -coverprofile=coverage.out && go tool cover -html=coverage.out

go build -o build/bin/webserver $(ls -1 cmd/web/*.go | grep -v _test.go)
./build/bin/webserver
