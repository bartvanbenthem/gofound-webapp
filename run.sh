#!/bin/bash

# cd tls
# go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

go test -vet=off -v ./cmd/web/
go run ./cmd/web/ --addr=":4000" \
                  --dsn="web:pass@/gofound?parseTime=true" \
                  --smtp-host="localhost" \
                  --smtp-port="1025" \
                  --smtp-user="" \
                  --smtp-password="" \
                  --mailaddr="mail@gofound.nl"