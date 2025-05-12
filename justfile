# Default task
default: help
API_FILE:= "./NoteVault.api"


# Setting
set dotenv-load

# Generate API documentation
gen-doc:
    @echo "Generating Swagger documentation..."
    goctl api swagger --api {{ API_FILE }} --dir api
    @echo "Generating default goctl API documentation..."
    goctl api doc --dir . -o api
    @echo "Documentation generation complete."

# Build the project with specified environment and drivers
build ENV="debug" DRIVERS="mysql,sqlite3,postgres":
    @echo "Building the project with environment: {{ENV}} and drivers: {{DRIVERS}}..."
    @if [ "{{ENV}}" = "debug" ]; then \
        go build -work -x -v -tags={{DRIVERS}} .; \
    elif [ "{{ENV}}" = "release" ]; then \
        go build -v -trimpath -tags={{DRIVERS}} -ldflags="-s -w" .; \
    else \
        echo "Invalid build environment. Use 'debug' or 'release'."; exit 1; \
    fi
    @echo "Build complete for drivers: {{DRIVERS}}."

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

# Help
[private]
help:
    @just --list

# hot-reload use air 
hot:
    @echo "Starting hot reload..."
    air -c .air.toml
    @echo "Hot reload started."

# Use goreleaser to build and release the project (snapshot)
snapshot:
    @echo "🛠 Build snapshot"
    goreleaser release --snapshot --clean
    @echo "Snapshot build complete."

# Release the project (If no tags are found, it will use snapshot)
goreleaser:
    @echo "🛠 Build release"
    @if [ -z "$(git tag --list)" ]; then \
        @echo "⚠️ No Git tags found, switching to snapshot release"; \
        just snapshot; \
    else \
        goreleaser release --clean; \
    fi
    @echo "Release build complete."

# Clean the go build binaries
clean:
    @echo "Cleaning up build artifacts, installed packages, and cache..."
    rm -rf dist/
    go clean -x
    @echo "Cleanup complete."

alias b := build
alias df := dep-fmt
alias g := gen-go
alias doc := gen-doc
alias n := snapshot
alias r := goreleaser
alias c := clean
alias i := init
alias h := hot
