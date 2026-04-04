.PHONY: build
build:
	@echo "Building pimpom..."
	@mkdir -p bin
	go build -o bin/pimpom main.go
	@echo "Build completed."

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...

.PHONY: build-all
build-all:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/pimpom-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/pimpom-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/pimpom-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/pimpom-darwin-arm64 main.go
	@echo "Builds complete in bin/"

.PHONY: format-imports
format-imports:
	GOFLAGS="-buildvcs=false" go run github.com/daixiang0/gci@latest write --skip-generated -s standard -s default -s "prefix(github.com/IvanJijon/pimpomodoro)" -s localmodule --custom-order .
