VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")

.PHONY: build
build:
	@echo "Building pimpom $(VERSION)..."
	@mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/pimpom main.go
	@echo "Build completed."

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: build-all
build-all:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/pimpom-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/pimpom-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/pimpom-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/pimpom-darwin-arm64 main.go
	@echo "Builds complete in bin/"

.PHONY: format-imports
format-imports:
	GOFLAGS="-buildvcs=false" go run github.com/daixiang0/gci@latest write --skip-generated -s standard -s default -s "prefix(github.com/IvanJijon/pimpomodoro)" -s localmodule --custom-order .

.PHONY: tag
tag:
	@if [ -z "$(V)" ]; then echo "Usage: make tag V=0.1.0"; exit 1; fi
	git tag -a v$(V) -m "Release v$(V)"
	@echo "Tagged v$(V)"

.PHONY: release
release: tag build-all
	@echo "Release v$(V) built. Push with: git push origin v$(V)"
