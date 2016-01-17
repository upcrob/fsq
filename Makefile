build:
	go tool yacc parser.y
	go build

install:
	go tool yacc parser.y
	go install

clean:
	rm -f test y.go y.output fsq
