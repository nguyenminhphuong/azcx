VERSION := 0.1.0
BINARY := azcx
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build install clean release test

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	cp $(BINARY) $(GOPATH)/bin/$(BINARY)

clean:
	rm -rf $(BINARY) dist/

test:
	go test -v ./...

# Cross-platform builds
release: clean
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe .

# Generate checksums
checksums:
	cd dist && shasum -a 256 * > checksums.txt

# Run locally
run:
	go run .

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run
