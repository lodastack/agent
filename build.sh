#/bin/bash 

version="0.1.2"
commit=`git rev-parse HEAD`
branch=`git rev-parse --abbrev-ref HEAD`
t=`date "+%Y-%m-%d_%H:%M:%S"`

cd cmd/agent

go build -v -ldflags="-X main.version=$version -X main.branch=$branch -X main.commit=$commit -X main.buildTime=$t"