# PowerShell build script for go-tenant-service
# Windows-compatible alternative to Makefile

param(
    [Parameter(Position=0)]
    [string]$Command = "help",
    
    [string]$DockerRegistry = "ghcr.io/vhvplatform",
    [string]$Version = ""
)

$ErrorActionPreference = "Stop"

$SERVICE_NAME = "tenant-service"
$GO_VERSION = "1.25.5"

if ($Version -eq "") {
    try {
        $Version = git describe --tags --always --dirty 2>$null
        if (-not $Version) {
            $Version = "dev"
        }
    } catch {
        $Version = "dev"
    }
}

function Show-Help {
    Write-Host ""
    Write-Host "go-tenant-service build script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\build.ps1 <command>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  help              Display this help message"
    Write-Host "  build             Build the service"
    Write-Host "  test              Run tests"
    Write-Host "  test-coverage     Run tests with coverage"
    Write-Host "  lint              Run linters (requires golangci-lint)"
    Write-Host "  fmt               Format code"
    Write-Host "  vet               Run go vet"
    Write-Host "  clean             Clean build artifacts"
    Write-Host "  run               Run the service locally"
    Write-Host "  deps              Download dependencies"
    Write-Host "  proto             Generate protobuf files"
    Write-Host "  docker-build      Build Docker image"
    Write-Host "  docker-run        Run Docker container locally"
    Write-Host "  install-tools     Install development tools"
    Write-Host ""
}

function Build-Service {
    Write-Host "Building $SERVICE_NAME..." -ForegroundColor Green
    
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    
    $outputPath = "bin\$SERVICE_NAME.exe"
    go build -o $outputPath .\cmd\main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build complete! Binary: $outputPath" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Run-Tests {
    Write-Host "Running tests..." -ForegroundColor Green
    go test -v -race ./...
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Tests failed!" -ForegroundColor Red
        exit 1
    }
}

function Run-TestsWithCoverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
    
    if ($LASTEXITCODE -eq 0) {
        go tool cover -html=coverage.txt -o coverage.html
        Write-Host "Coverage report generated: coverage.html" -ForegroundColor Green
    } else {
        Write-Host "Tests failed!" -ForegroundColor Red
        exit 1
    }
}

function Run-Lint {
    Write-Host "Running linters..." -ForegroundColor Green
    
    $golangciLint = Get-Command golangci-lint -ErrorAction SilentlyContinue
    if (-not $golangciLint) {
        Write-Host "golangci-lint not found. Please install it first using: .\build.ps1 install-tools" -ForegroundColor Yellow
        exit 1
    }
    
    golangci-lint run ./...
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Linting failed!" -ForegroundColor Red
        exit 1
    }
}

function Format-Code {
    Write-Host "Formatting code..." -ForegroundColor Green
    go fmt ./...
    gofmt -s -w .
    Write-Host "Formatting complete!" -ForegroundColor Green
}

function Run-Vet {
    Write-Host "Running go vet..." -ForegroundColor Green
    go vet ./...
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "go vet failed!" -ForegroundColor Red
        exit 1
    }
}

function Clean-Build {
    Write-Host "Cleaning..." -ForegroundColor Green
    
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force bin
    }
    if (Test-Path "dist") {
        Remove-Item -Recurse -Force dist
    }
    if (Test-Path "coverage.txt") {
        Remove-Item coverage.txt
    }
    if (Test-Path "coverage.html") {
        Remove-Item coverage.html
    }
    
    $outFiles = Get-ChildItem -Filter "*.out" -File
    foreach ($file in $outFiles) {
        Remove-Item $file.FullName
    }
    
    go clean -testcache
    Write-Host "Clean complete!" -ForegroundColor Green
}

function Run-Service {
    Write-Host "Running $SERVICE_NAME..." -ForegroundColor Green
    go run .\cmd\main.go
}

function Download-Dependencies {
    Write-Host "Downloading dependencies..." -ForegroundColor Green
    go mod download
    go mod tidy
    Write-Host "Dependencies downloaded!" -ForegroundColor Green
}

function Generate-Proto {
    if (Test-Path "proto") {
        Write-Host "Generating protobuf files..." -ForegroundColor Green
        
        $protocCmd = Get-Command protoc -ErrorAction SilentlyContinue
        if (-not $protocCmd) {
            Write-Host "protoc not found. Please install Protocol Buffers compiler." -ForegroundColor Yellow
            exit 1
        }
        
        $protoFiles = Get-ChildItem -Path "proto" -Filter "*.proto"
        foreach ($file in $protoFiles) {
            protoc --go_out=. --go_opt=paths=source_relative `
                --go-grpc_out=. --go-grpc_opt=paths=source_relative `
                $file.FullName
        }
        
        Write-Host "Protobuf generation complete!" -ForegroundColor Green
    } else {
        Write-Host "No proto directory found, skipping..." -ForegroundColor Yellow
    }
}

function Build-Docker {
    Write-Host "Building Docker image..." -ForegroundColor Green
    
    $dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
    if (-not $dockerCmd) {
        Write-Host "Docker not found. Please install Docker Desktop for Windows." -ForegroundColor Yellow
        exit 1
    }
    
    $imageTag = "$DockerRegistry/${SERVICE_NAME}:$Version"
    docker build -t $imageTag .
    docker tag $imageTag "${DockerRegistry}/${SERVICE_NAME}:latest"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Docker image built: $imageTag" -ForegroundColor Green
    } else {
        Write-Host "Docker build failed!" -ForegroundColor Red
        exit 1
    }
}

function Run-Docker {
    Write-Host "Running Docker container..." -ForegroundColor Green
    
    $dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
    if (-not $dockerCmd) {
        Write-Host "Docker not found. Please install Docker Desktop for Windows." -ForegroundColor Yellow
        exit 1
    }
    
    docker run --rm -p 8080:8080 -p 50051:50051 `
        --name $SERVICE_NAME `
        "${DockerRegistry}/${SERVICE_NAME}:latest"
}

function Install-Tools {
    Write-Host "Installing development tools..." -ForegroundColor Green
    
    Write-Host "Installing golangci-lint..." -ForegroundColor Yellow
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    if (Test-Path "proto") {
        Write-Host "Installing protobuf tools..." -ForegroundColor Yellow
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    }
    
    Write-Host "Tools installed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Note: Make sure your GOPATH\bin is in your PATH environment variable" -ForegroundColor Yellow
    Write-Host "Typically: C:\Users\<YourUsername>\go\bin" -ForegroundColor Yellow
}

# Main script execution
switch ($Command.ToLower()) {
    "help"            { Show-Help }
    "build"           { Build-Service }
    "test"            { Run-Tests }
    "test-coverage"   { Run-TestsWithCoverage }
    "lint"            { Run-Lint }
    "fmt"             { Format-Code }
    "vet"             { Run-Vet }
    "clean"           { Clean-Build }
    "run"             { Run-Service }
    "deps"            { Download-Dependencies }
    "proto"           { Generate-Proto }
    "docker-build"    { Build-Docker }
    "docker-run"      { Run-Docker }
    "install-tools"   { Install-Tools }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}
