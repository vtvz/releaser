build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/releaser-darwin-amd64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/releaser-linux-amd64

build: build-linux build-mac
