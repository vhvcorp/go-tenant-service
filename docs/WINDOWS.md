# Windows Development Guide

This guide provides instructions for developing the `go-tenant-service` on Windows.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Building](#building)
- [Running the Service](#running-the-service)
- [Testing](#testing)
- [Development Tools](#development-tools)
- [Docker on Windows](#docker-on-windows)
- [Common Issues](#common-issues)
- [IDE Setup](#ide-setup)

## Prerequisites

### Required Software

1. **Go 1.25.5 or later**
   - Download from: https://go.dev/dl/
   - During installation, ensure "Add to PATH" is checked
   - Verify installation: `go version`

2. **Git for Windows**
   - Download from: https://git-scm.com/download/win
   - Recommended: Use Git Bash for Unix-like command experience
   - Verify installation: `git --version`

3. **PowerShell 5.1 or later** (Pre-installed on Windows 10/11)
   - Check version: `$PSVersionTable.PSVersion`
   - Optional: Install PowerShell 7+ for better experience

### Optional Software

4. **Docker Desktop for Windows** (for containerization)
   - Download from: https://www.docker.com/products/docker-desktop
   - Requires Windows 10/11 Pro, Enterprise, or Education
   - WSL 2 backend recommended

5. **MongoDB** (if running locally)
   - Download from: https://www.mongodb.com/try/download/community
   - Or use Docker: `docker run -d -p 27017:27017 mongo:7.0`

6. **Redis** (if running locally)
   - Download from: https://github.com/microsoftarchive/redis/releases
   - Or use Docker: `docker run -d -p 6379:6379 redis:7-alpine`

7. **Visual Studio Code** (recommended IDE)
   - Download from: https://code.visualstudio.com/
   - Install Go extension: `ms-vscode.go`

## Installation

### 1. Clone the Repository

Using PowerShell or Git Bash:

```powershell
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service
```

### 2. Download Dependencies

Using PowerShell:

```powershell
.\build.ps1 deps
```

Or using Go directly:

```powershell
go mod download
go mod tidy
```

### 3. Install Development Tools

```powershell
.\build.ps1 install-tools
```

This installs:
- `golangci-lint` for code linting
- `protoc-gen-go` and `protoc-gen-go-grpc` for protobuf (if needed)

**Important:** Ensure your Go bin directory is in your PATH:
- Typically: `C:\Users\<YourUsername>\go\bin`
- Add to PATH via System Environment Variables or in PowerShell:

```powershell
$env:PATH += ";$env:USERPROFILE\go\bin"
```

To make it permanent, add to your PowerShell profile:

```powershell
notepad $PROFILE
# Add the line: $env:PATH += ";$env:USERPROFILE\go\bin"
```

## Building

### Using PowerShell Script

```powershell
# Build the service
.\build.ps1 build

# The binary will be created at: .\bin\tenant-service.exe
```

### Using Go Directly

```powershell
go build -o bin\tenant-service.exe .\cmd\main.go
```

### View All Build Commands

```powershell
.\build.ps1 help
```

### Using Batch File (cmd.exe)

If you prefer using Command Prompt:

```batch
build.bat build
build.bat test
build.bat help
```

## Running the Service

### Prerequisites

Before running the service, ensure the required services are available:

1. **MongoDB** (default port: 27017)
2. **Redis** (default port: 6379)

### Set Environment Variables

Create a `.env` file or set environment variables in PowerShell:

```powershell
# Server Configuration
$env:TENANT_SERVICE_PORT = "50053"
$env:TENANT_SERVICE_HTTP_PORT = "8083"

# Database Configuration
$env:MONGODB_URI = "mongodb://localhost:27017"
$env:MONGODB_DATABASE = "saas_framework"

# Redis Configuration
$env:REDIS_URL = "redis://localhost:6379/0"

# Logging
$env:LOG_LEVEL = "info"
```

### Run the Service

Using the build script:

```powershell
.\build.ps1 run
```

Or run directly:

```powershell
go run .\cmd\main.go
```

Or run the compiled binary:

```powershell
.\bin\tenant-service.exe
```

### Verify the Service is Running

Open a new PowerShell window and test the health endpoint:

```powershell
# HTTP Health Check
Invoke-WebRequest -Uri http://localhost:8083/health

# Or using curl (if installed via Git for Windows)
curl http://localhost:8083/health
```

## Testing

### Run All Tests

```powershell
.\build.ps1 test
```

### Run Tests with Coverage

```powershell
.\build.ps1 test-coverage
```

This generates `coverage.html` that you can open in a browser.

### Run Tests with Go

```powershell
go test -v ./...
go test -v -race ./...
```

## Development Tools

### Code Formatting

```powershell
# Format all code
.\build.ps1 fmt
```

### Linting

```powershell
# Run linters
.\build.ps1 lint
```

### Go Vet

```powershell
# Run go vet
.\build.ps1 vet
```

### Clean Build Artifacts

```powershell
# Clean all build artifacts
.\build.ps1 clean
```

## Docker on Windows

### Prerequisites

- Docker Desktop for Windows installed and running
- WSL 2 backend enabled (recommended)

### Build Docker Image

```powershell
.\build.ps1 docker-build
```

Or using Docker directly:

```powershell
docker build -t ghcr.io/vhvplatform/tenant-service:latest .
```

### Run Docker Container

```powershell
.\build.ps1 docker-run
```

Or using Docker directly:

```powershell
docker run --rm -p 8080:8080 -p 50051:50051 `
  --name tenant-service `
  ghcr.io/vhvplatform/tenant-service:latest
```

### Using Docker Compose (if available)

```powershell
docker-compose up -d
```

## Common Issues

### Issue: "go: command not found"

**Solution:** Go is not in your PATH. Reinstall Go and ensure "Add to PATH" is checked, or manually add Go to PATH:

```powershell
$env:PATH += ";C:\Program Files\Go\bin"
```

### Issue: "execution of scripts is disabled on this system"

**Solution:** PowerShell execution policy is too restrictive. Run PowerShell as Administrator:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

Or run scripts with bypass:

```powershell
powershell -ExecutionPolicy Bypass -File .\build.ps1 build
```

### Issue: "Cannot connect to MongoDB"

**Solution:** Ensure MongoDB is running:

```powershell
# Check if MongoDB is running
Get-Service -Name MongoDB

# Or start MongoDB service
Start-Service MongoDB
```

If using Docker:

```powershell
docker run -d -p 27017:27017 --name mongodb mongo:7.0
```

### Issue: Line Ending Issues (CRLF vs LF)

**Solution:** Git may convert line endings. Configure Git to handle line endings:

```powershell
git config --global core.autocrlf true
```

For this repository:

```powershell
git config core.autocrlf input
```

### Issue: "golangci-lint: command not found"

**Solution:** The Go bin directory is not in your PATH. Add it:

```powershell
$env:PATH += ";$env:USERPROFILE\go\bin"
```

### Issue: Port Already in Use

**Solution:** Another process is using the port. Find and stop it:

```powershell
# Find process using port 8083
netstat -ano | findstr :8083

# Kill the process (replace PID with actual process ID)
Stop-Process -Id <PID> -Force
```

### Issue: File Paths with Backslashes

**Solution:** Go handles both forward slashes and backslashes on Windows. Use forward slashes in code for cross-platform compatibility:

```go
// Good (cross-platform)
path := "internal/service/tenant.go"

// Also works on Windows
path := "internal\\service\\tenant.go"
```

## IDE Setup

### Visual Studio Code

1. **Install Go Extension**
   ```powershell
   code --install-extension golang.go
   ```

2. **Recommended Settings** (`.vscode/settings.json`):
   ```json
   {
     "go.useLanguageServer": true,
     "go.lintTool": "golangci-lint",
     "go.lintOnSave": "package",
     "go.formatTool": "gofmt",
     "editor.formatOnSave": true,
     "[go]": {
       "editor.codeActionsOnSave": {
         "source.organizeImports": true
       }
     }
   }
   ```

3. **Install Go Tools**
   - Open Command Palette: `Ctrl+Shift+P`
   - Type: "Go: Install/Update Tools"
   - Select all tools and install

### GoLand / IntelliJ IDEA

1. Install Go plugin (if using IntelliJ IDEA)
2. Open project directory
3. Configure Go SDK: `File > Settings > Go > GOROOT`
4. Enable Go modules: `File > Settings > Go > Go Modules`

### Debugging

#### VS Code Launch Configuration (`.vscode/launch.json`)

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Service",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/main.go",
      "env": {
        "TENANT_SERVICE_PORT": "50053",
        "TENANT_SERVICE_HTTP_PORT": "8083",
        "MONGODB_URI": "mongodb://localhost:27017",
        "MONGODB_DATABASE": "saas_framework",
        "REDIS_URL": "redis://localhost:6379/0",
        "LOG_LEVEL": "debug"
      }
    }
  ]
}
```

## Performance Considerations

### Windows-Specific Notes

1. **File System Performance:** Windows file system operations can be slower than Unix-based systems. Consider:
   - Using SSD for development
   - Disabling real-time antivirus scanning for the project directory
   - Using WSL 2 for better performance (optional)

2. **Docker Performance:** 
   - Use WSL 2 backend for better Docker performance
   - Store projects in WSL 2 file system for faster builds

3. **Build Times:**
   - First build may be slower due to dependency downloads
   - Subsequent builds use cached modules

## Additional Resources

- [Go on Windows Documentation](https://go.dev/doc/install/windows)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/install/)
- [PowerShell Documentation](https://docs.microsoft.com/en-us/powershell/)
- [VS Code Go Extension](https://code.visualstudio.com/docs/languages/go)

## Getting Help

If you encounter issues not covered in this guide:

1. Check the main [README.md](../README.md)
2. Review [CONTRIBUTING.md](../CONTRIBUTING.md)
3. Search existing [GitHub Issues](https://github.com/vhvplatform/go-tenant-service/issues)
4. Open a new issue with:
   - Windows version
   - Go version (`go version`)
   - PowerShell version (`$PSVersionTable.PSVersion`)
   - Error messages and logs

## Contributing

When contributing from Windows:

1. Ensure code passes linting: `.\build.ps1 lint`
2. Format code properly: `.\build.ps1 fmt`
3. Run tests: `.\build.ps1 test`
4. Follow cross-platform best practices:
   - Use `filepath.Join()` for paths
   - Avoid platform-specific system calls
   - Test on multiple platforms if possible
