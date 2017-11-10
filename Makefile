name=fsq
version=1.7.1

build: genparser
	go build

install: genparser
	go install

all: test
	rm -Rf bin
	mkdir bin
	export GOOS=darwin; export GOARCH=amd64; go build -o bin/$(name)-$(version)-osx-amd64
	export GOOS=linux; export GOARCH=386; go build -o bin/$(name)-$(version)-linux-386
	export GOOS=linux; export GOARCH=amd64; go build -o bin/$(name)-$(version)-linux-amd64
	export GOOS=linux; export GOARCH=arm; go build -o bin/$(name)-$(version)-linux-arm
	export GOOS=linux; export GOARCH=arm64; go build -o bin/$(name)-$(version)-linux-arm64
	export GOOS=windows; export GOARCH=386; go build -o bin/$(name)-$(version)-windows-386.exe
	export GOOS=windows; export GOARCH=amd64; go build -o bin/$(name)-$(version)-windows-amd64.exe
	export GOOS=freebsd; export GOARCH=amd64; go build -o bin/$(name)-$(version)-freebsd-amd64
	export GOOS=freebsd; export GOARCH=386; go build -o bin/$(name)-$(version)-freebsd-386

test: genparser
	go test

genparser: installyacc
	goyacc parser.y

installyacc:
	go get golang.org/x/tools/cmd/goyacc
	go install golang.org/x/tools/cmd/goyacc

clean:
	rm -f y.go y.output fsq
	rm -Rf bin
