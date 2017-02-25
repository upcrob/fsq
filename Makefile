name=fsq
version=1.3.0

build: test
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

install: installyacc
	goyacc parser.y
	go install

test: installyacc
	goyacc parser.y
	go test

installyacc:
	go get golang.org/x/tools/cmd/goyacc
	go install golang.org/x/tools/cmd/goyacc

clean:
	rm -f y.go y.output
	rm -Rf bin
