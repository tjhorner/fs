.PHONY: dist dist-win dist-macos dist-linux ensure-dist-dir build install uninstall

GOBUILD=go build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

# Change me!
PROJECT_NAME=fs

ensure-dist-dir:
	@- mkdir -p dist

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-windows-amd64.exe *.go

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-darwin-amd64 *.go

dist-linux: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-linux-amd64 *.go

dist: dist-win dist-macos dist-linux

build:
	@- mkdir -p bin
	$(GOBUILD) -o bin/$(PROJECT_NAME) *.go
	@- chmod +x bin/$(PROJECT_NAME)

install: build
	mv bin/$(PROJECT_NAME) $(INSTALLPATH)/$(PROJECT_NAME)
	@- rm -rf bin
	@echo "$(PROJECT_NAME) was installed to $(INSTALLPATH)/$(PROJECT_NAME). Run make uninstall to get rid of it, or just remove the binary yourself."

uninstall:
	rm $(INSTALLPATH)/$(PROJECT_NAME)

run:
	@- go run *.go