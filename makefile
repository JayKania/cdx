.PHONY: build-linux build-darwin build-windows build-all run-linux run-darwin run-windows clean test fmt

build-darwin:
	rm -rf build/darwin/
	go build -o build/darwin/cdx ./cmd/.
	go build -o build/darwin/setup ./setup/.

build-linux:
	rm -rf build/linux/
	go build -o build/linux/cdx ./cmd/.
	go build -o build/linux/setup ./setup/.

build-windows:
	rm -rf build/windows/
	

	go build -o build/windows/cdx ./cmd/.


	go build -o build/windows/setup ./setup/.

build-all:
	rm -rf build/


	GOOS=linux GOARCH=amd64 go build -o build/linux/cdx ./cmd/.
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/cdx ./cmd/.
	GOOS=windows GOARCH=amd64 go build -o build/windows/cdx ./cmd/.


	GOOS=linux GOARCH=amd64 go build -o build/linux/setup ./setup/.
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/setup ./setup/.
	GOOS=windows GOARCH=amd64 go build -o build/windows/setup ./setup/.

run-linux:

	GOOS=linux GOARCH=amd64 go run ./cmd/.

run-darwin:

	GOOS=darwin GOARCH=amd64 go run ./cmd/.

run-windows:

	GOOS=windows GOARCH=amd64 go run ./cmd/.

clean:
	rm -rf ./cmd/build/

test:
	go test ./...

fmt:
	go fmt *.go