all: build

fmt:
	gofmt -l -w -s ./

dep:fmt
	go get github.com/oiooj/cli
	go get github.com/lodastack/nux
	go get github.com/lodastack/log
	go get github.com/lodastack/models
	go get github.com/BurntSushi/toml

install:dep
	go install agent

build:fmt
	cd cmd/agent && go build -v

clean:
	cd cmd/agent && go clean
