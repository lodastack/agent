#!/bin/bash

go version

# install dep

go get github.com/golang/dep && cd ~/.go_workspace/src/github.com/golang/dep/cmd/dep && go build -v && sudo cp dep /usr/local/bin/
cd ${HOME}/.go_workspace/src/github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME && dep ensure -v

go test -timeout 60s -v ./...

# build
cd ./cmd/agent/
go build -v
