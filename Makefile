.PHONY: build
build:
	@echo "Building pimpom..."
	go build -o pimpom main.go
	@echo "Build completed."

.PHONY: run
run:
	@echo "Running pimpom..."
	go run main.go
	@echo "Run completed."

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -o pimpom-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o pimpom-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o pimpom-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o pimpom-darwin-arm64 main.go
	@echo "Builds complete in current directory"

.PHONY: format-imports
format-imports:
	GOFLAGS="-buildvcs=false" go run github.com/daixiang0/gci@latest write --skip-generated -s standard -s default -s "prefix(github.com/IvanJijon/pimpomodoro)" -s localmodule --custom-order .
