name=fsq

build:
	go tool yacc parser.y

	rm -Rf bin
	mkdir bin

	export GOOS=darwin; export GOARCH=amd64; go build -o bin/$(name)-osx-amd64
	export GOOS=linux; export GOARCH=386; go build -o bin/$(name)-linux-386
	export GOOS=linux; export GOARCH=amd64; go build -o bin/$(name)-linux-amd64
	export GOOS=linux; export GOARCH=arm; go build -o bin/$(name)-linux-arm
	export GOOS=linux; export GOARCH=arm64; go build -o bin/$(name)-linux-arm64
	export GOOS=windows; export GOARCH=386; go build -o bin/$(name)-windows-386.exe
	export GOOS=windows; export GOARCH=amd64; go build -o bin/$(name)-windows-amd64.exe
	export GOOS=freebsd; export GOARCH=amd64; go build -o bin/$(name)-freebsd-amd64
	export GOOS=freebsd; export GOARCH=386; go build -o bin/$(name)-freebsd-386

install:
	go tool yacc parser.y
	go install

clean:
	rm -f y.go y.output
	rm -Rf bin
