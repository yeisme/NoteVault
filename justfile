# Default task
[private]
default: help
API_FILE:= "./notevault.api"


# Setting
set dotenv-load

# Generate all artifacts: Generate documentation and Golang code
gen-all: gen-doc gen-go gen-sql2go

# Generate API documentation
gen-doc:
    @echo "Generating Swagger documentation..."
    rm api/notevault.*
    goctl api swagger --api {{ API_FILE }} --dir api
    @echo "Generating default goctl API documentation..."
    goctl api doc --dir . -o api
    @echo "Documentation generation complete."

# Build the project with specified environment and external drivers (e.g., mysql, sqlite3)
build ENV="debug" DRIVERS="mysql": dep-fmt
    @echo "Building the project with environment: {{ENV}} and external drivers: {{DRIVERS}}..."
    @if [ "{{ENV}}" = "debug" ]; then \
        go build -work -x -v -tags={{DRIVERS}} .; \
    elif [ "{{ENV}}" = "release" ]; then \
        go build -v -trimpath -tags={{DRIVERS}} -ldflags="-s -w" .; \
    else \
        echo "Invalid build environment. Use 'debug' or 'release'."; exit 1; \
    fi
    @echo "Build complete for drivers: {{DRIVERS}}."

# Alias for building with sqlite3 driver
b3:
    just build debug sqlite3

# Update dependencies and tidy up the package
dep-fmt:
    @echo "Formatting code..."
    gofmt -w .
    @echo "Get dependencies..."
    go mod tidy
    @echo "Code formatting and dependency update complete."

# Download dependencies and tools
init:
    @echo "Downloading goctl..."
    go install github.com/zeromicro/go-zero/tools/goctl@latest
    @echo "Downloading gorm/gen..."
    go install gorm.io/gen/tools/gentool@latest
    @echo "Install air..."
    go install github.com/air-verse/air@latest
    @echo "Install goreleaser..."
    go install github.com/goreleaser/goreleaser/v2@latest
    @echo "Initialization complete."

# Generate Golang code from the API file
gen-go:
    @echo "Generating API file..."
    goctl api go --api {{ API_FILE }} --dir .
    @echo "API file generation complete."

# Generate Golang code from the sql file to the gorm gen file
gen-sql2go:
    go generate ./...
    @echo "SQL to Go generation complete."

# Help
help:
    @just --list

# hot-reload use air 
hot: dep-fmt
    @echo "Starting hot reload..."
    air -c .air.toml
    @echo "Hot reload started."

# Use goreleaser to build and release the project (snapshot)
snapshot: dep-fmt
    @echo "🛠 Build snapshot"
    goreleaser release --snapshot --clean
    @echo "Snapshot build complete."

# Release the project (local, remote, or test)
goreleaser MODE="local": dep-fmt
    @echo "🛠 Build release"
    @if [ "{{MODE}}" = "local" ]; then \
        just snapshot; \
    elif [ "{{MODE}}" = "remote" ]; then \
        goreleaser release --clean; \
    elif [ "{{MODE}}" = "test" ]; then \
        goreleaser release --clean --skip=publish; \
    else \
        echo "Invalid mode. Use 'local' or 'remote'."; exit 1; \
    fi
    @echo "Release build complete."

# Clean the go build binaries
clean:
    @echo "Cleaning up build artifacts, installed packages, and cache..."
    rm -rf dist/
    rm -rf tmp/
    go clean -x
    @echo "Cleaning up logs..."
    rm -rf logs/*
    @echo "Cleanup complete."

# Validate configuration files
validate:
    @echo "Validating configuration files..."
    go run . validate
    @echo "Configuration validation complete."

# Set a Development environment with docker compose
docker-dev:
    @echo "Setting up development environment with Docker Compose..."
    docker-compose -f script/dev/docker-compose.yaml up -d
    @echo "Development environment is up and running."

alias dev := docker-dev
alias b := build
alias df := dep-fmt
alias g := gen-go
alias doc := gen-doc
alias n := snapshot
alias r := goreleaser
alias c := clean
alias i := init
alias h := hot
alias v := validate
alias ga := gen-all
alias g2 := gen-sql2go
