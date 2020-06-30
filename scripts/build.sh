#!/bin/sh

# compile the runners with gcc disabled
export CGO_ENABLED=0
go build -o release/${GOOS}/${GOARCH}/go-wol-server
