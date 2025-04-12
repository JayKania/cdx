.PHONY: build build-all clean run test

build-darwin:
	rm -rf build/darwin/
	go build -o build/darwin/cdx .

build-linux:
	rm -rf build/linux/
	go build -o build/linux/cdx .

build-windows:
	rm -rf build/windows/
	go build -o build/windows/cdx .

build-all:
	rm -rf build/
	GOOS=linux GOARCH=amd64 go build -o build/linux/cdx .
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/cdx .
	GOOS=windows GOARCH=amd64 go build -o build/windows/cdx .

run-linux:
	GOOS=linux GOARCH=amd64 go run .

run-darwin:
	GOOS=darwin GOARCH=amd64 go run .

run-windows:
	GOOS=windows GOARCH=amd64 go run .

clean:
	rm -rf build/

test:
	go test ./...

fmt:
	go fmt *.go