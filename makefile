.PHONY: build build-all clean run test

build:
	go build -o build/cdx main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o build/cdx-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o build/cdx-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o build/cdx-windows-amd64.exe main.go

clean:
	rm -rf build/

run:
	go run main.go

test:
	go test ./...