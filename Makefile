build-linux:
	GOOS=darwin GOARCH=amd64 go build -o bin/releaser-darwin-amd64

build-mac:
	GOOS=linux GOARCH=amd64 go build -o bin/releaser-linux-amd64

build: build-mac build-linux