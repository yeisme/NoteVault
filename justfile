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

# Placeholder for build task, assuming you might have one
build ENV="debug": dep-fmt
    @echo "Building the project ({{ENV}})..."
    @if [ "{{ENV}}" = "debug" ]; then \
        go build -work -x -v .; \
    elif [ "{{ENV}}" = "release" ]; then \
        go build -v -trimpath -ldflags="-s -w" .; \
    else \
        echo "Invalid build environment. Use 'debug' or 'release'."; exit 1; \
    fi
    @echo "Build complete."

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
help:
    @just --list

# hot-reload use air 
hot:
    @echo "Starting hot reload..."
    air -c .air.toml
    @echo "Hot reload started."

# Use goreleaser to build and release the project
snapshot:
    @echo "🛠 Build snapshot"
    goreleaser release --snapshot --clean
    @echo "Snapshot build complete."

# Release the project（若无 Git 标签则 fallback 到 snapshot）
release:
    @echo "🛠 Build release"
    @if [ -z "$(git tag --list)" ]; then \
        echo "⚠️ No Git tags found, switching to snapshot release"; \
        just snapshot; \
    else \
        goreleaser release --clean; \
    fi
    @echo "Release build complete."

# Clean the go build binaries
clean:
    @echo "Cleaning up build artifacts, installed packages, and cache..."
    go clean -x
    @echo "Cleanup complete."

alias b := build
alias d := dep-fmt
alias g := gen-go
alias doc := gen-doc
alias n := snapshot
alias r := release
alias c := clean
alias i := init
alias h := hot
