.PHONY: build test clean

VERSION=$(shell git describe --tags)

build: terraless_amd64 terraless_alpine terraless_darwin terraless.exe

.get-deps: *.go
	go get -t -d -v ./...
	touch .get-deps

clean:
	rm -f .get-deps
	rm -f *_amd64 *_darwin *.exe

test: .get-deps *.go
	go test -v ./...

terraless_amd64: .get-deps *.go
	 GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $@ *.go

terraless_alpine: .get-deps *.go
	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ *.go

terraless_darwin: .get-deps *.go
	GOOS=darwin go build -o $@ *.go

terraless.exe: .get-deps *.go
	GOOS=windows GOARCH=amd64 go build -o $@ *.go
