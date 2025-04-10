.PHONY: build build-all clean run test

build:
	go build -o build/cdx *.go

build-home:
	rm -f ~/cdx
	go build -o ~/cdx *.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o build/cdx-linux-amd64 *.go
	GOOS=darwin GOARCH=amd64 go build -o build/cdx-darwin-amd64 *.go
	GOOS=windows GOARCH=amd64 go build -o build/cdx-windows-amd64.exe *.go

clean:
	rm -rf build/

run:
	go run *.go

test:
	go test ./...

fmt:
	go fmt *.go