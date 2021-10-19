#!/bin/bash

# go test -v -cover cmd/web/*
# go test -v -cover internal/handlers/*
# go test -v -cover internal/render/*
# go test -v -cover internal/forms/*

# go test -coverprofile=coverage.out && go tool cover -html=coverage.out

go build -o build/bin/webserver $(ls -1 cmd/web/*.go | grep -v _test.go)
./build/bin/webserver
