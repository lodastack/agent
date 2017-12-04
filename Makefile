all: build

test:
	go test -timeout 60s -v ./...

fmt:
	gofmt -l -w -s ./

dep:fmt
	gdm restore

install:dep
	go install agent

build:dep
	./build.sh

clean:
	cd cmd/agent && go clean
