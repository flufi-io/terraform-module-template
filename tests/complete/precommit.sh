#!/bin/bash
cd tests/complete || exit
rm -rf go.*
export GOOS=linux GOARCH=arm64 CGO_ENABLED=0
go mod init test
go mod tidy >> /dev/null
go test -count=1 -timeout 120m -v ./...
exit $?
