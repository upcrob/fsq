name=fsq
version=1.8.1

build: genparser
	go build

install: genparser
	go install

all: test
	rm -Rf bin
	mkdir -p bin/$(name)-$(version)-osx-amd64
	mkdir -p bin/$(name)-$(version)-linux-386
	mkdir -p bin/$(name)-$(version)-linux-amd64
	mkdir -p bin/$(name)-$(version)-linux-arm
	mkdir -p bin/$(name)-$(version)-linux-arm64
	mkdir -p bin/$(name)-$(version)-windows-386
	mkdir -p bin/$(name)-$(version)-windows-amd64
	mkdir -p bin/$(name)-$(version)-freebsd-amd64
	mkdir -p bin/$(name)-$(version)-freebsd-386

	export GOOS=darwin; export GOARCH=amd64; go build -o bin/$(name)-$(version)-osx-amd64/$(name)
	export GOOS=linux; export GOARCH=386; go build -o bin/$(name)-$(version)-linux-386/$(name)
	export GOOS=linux; export GOARCH=amd64; go build -o bin/$(name)-$(version)-linux-amd64/$(name)
	export GOOS=linux; export GOARCH=arm; go build -o bin/$(name)-$(version)-linux-arm/$(name)
	export GOOS=linux; export GOARCH=arm64; go build -o bin/$(name)-$(version)-linux-arm64/$(name)
	export GOOS=windows; export GOARCH=386; go build -o bin/$(name)-$(version)-windows-386/$(name).exe
	export GOOS=windows; export GOARCH=amd64; go build -o bin/$(name)-$(version)-windows-amd64/$(name).exe
	export GOOS=freebsd; export GOARCH=amd64; go build -o bin/$(name)-$(version)-freebsd-amd64/$(name)
	export GOOS=freebsd; export GOARCH=386; go build -o bin/$(name)-$(version)-freebsd-386/$(name)

	(cd bin && zip -r - $(name)-$(version)-osx-amd64) > bin/$(name)-$(version)-osx-amd64.zip
	(cd bin && zip -r - $(name)-$(version)-linux-386) > bin/$(name)-$(version)-linux-386.zip
	(cd bin && zip -r - $(name)-$(version)-linux-amd64) > bin/$(name)-$(version)-linux-amd64.zip
	(cd bin && zip -r - $(name)-$(version)-linux-arm) > bin/$(name)-$(version)-linux-arm.zip
	(cd bin && zip -r - $(name)-$(version)-linux-arm64) > bin/$(name)-$(version)-linux-arm64.zip
	(cd bin && zip -r - $(name)-$(version)-windows-386) > bin/$(name)-$(version)-windows-386.zip
	(cd bin && zip -r - $(name)-$(version)-windows-amd64) > bin/$(name)-$(version)-windows-amd64.zip
	(cd bin && zip -r - $(name)-$(version)-freebsd-amd64) > bin/$(name)-$(version)-freebsd-amd64.zip
	(cd bin && zip -r - $(name)-$(version)-freebsd-386) > bin/$(name)-$(version)-freebsd-386.zip

	rm -Rf bin/$(name)-$(version)-osx-amd64
	rm -Rf bin/$(name)-$(version)-linux-386
	rm -Rf bin/$(name)-$(version)-linux-amd64
	rm -Rf bin/$(name)-$(version)-linux-arm
	rm -Rf bin/$(name)-$(version)-linux-arm64
	rm -Rf bin/$(name)-$(version)-windows-386
	rm -Rf bin/$(name)-$(version)-windows-amd64
	rm -Rf bin/$(name)-$(version)-freebsd-amd64
	rm -Rf bin/$(name)-$(version)-freebsd-386

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
