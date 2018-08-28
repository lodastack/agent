#!/bin/bash

export GO111MODULE=on
go mod download
go mod verify
go test -timeout 60s -v ./...
# build
cd ./cmd/agent/ && go build -v -mod readonly
