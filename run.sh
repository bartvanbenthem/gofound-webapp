#!/bin/bash

go test -vet=off -v ./cmd/web/
go run ./cmd/web/ --addr=":4000" \
                  --dsn="web:pass@/gofound?parseTime=true" \
                  --smtp-host="localhost" \
                  --smtp-port="1025" \
                  --smtp-user="" \
                  --smtp-password="" \
                  --mail-address="mail@gofound.nl" \
                  --cert="./tls/cert.pem" \
                  --key="./tls/key.pem"
