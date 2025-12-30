# Windows Compatibility Quick Reference

This document provides a quick reference for common tasks on Windows.

## Quick Start

```powershell
# Clone and setup
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service
.\build.ps1 deps

# Build
.\build.ps1 build

# Run
.\build.ps1 run
```

## Command Comparison

| Task | Linux/macOS | Windows PowerShell | Windows CMD |
|------|-------------|-------------------|-------------|
| Build | `make build` | `.\build.ps1 build` | `build.bat build` |
| Test | `make test` | `.\build.ps1 test` | `build.bat test` |
| Run | `make run` | `.\build.ps1 run` | `build.bat run` |
| Lint | `make lint` | `.\build.ps1 lint` | `build.bat lint` |
| Format | `make fmt` | `.\build.ps1 fmt` | `build.bat fmt` |
| Clean | `make clean` | `.\build.ps1 clean` | `build.bat clean` |
| Docker Build | `make docker-build` | `.\build.ps1 docker-build` | `build.bat docker-build` |
| Docker Run | `make docker-run` | `.\build.ps1 docker-run` | `build.bat docker-run` |
| Install Tools | `make install-tools` | `.\build.ps1 install-tools` | `build.bat install-tools` |

## Build Script Commands

```powershell
.\build.ps1 help              # Show all available commands
.\build.ps1 build             # Build the service
.\build.ps1 test              # Run tests
.\build.ps1 test-coverage     # Run tests with coverage report
.\build.ps1 lint              # Run golangci-lint
.\build.ps1 fmt               # Format code with gofmt
.\build.ps1 vet               # Run go vet
.\build.ps1 clean             # Clean build artifacts
.\build.ps1 run               # Run the service
.\build.ps1 deps              # Download dependencies
.\build.ps1 proto             # Generate protobuf files
.\build.ps1 docker-build      # Build Docker image
.\build.ps1 docker-run        # Run Docker container
.\build.ps1 install-tools     # Install development tools
```

## Environment Variables

### PowerShell

```powershell
# Set for current session
$env:TENANT_SERVICE_PORT = "50053"
$env:MONGODB_URI = "mongodb://localhost:27017"

# Set permanently (user level)
[System.Environment]::SetEnvironmentVariable("TENANT_SERVICE_PORT", "50053", "User")
```

### CMD

```batch
REM Set for current session
set TENANT_SERVICE_PORT=50053
set MONGODB_URI=mongodb://localhost:27017

REM Set permanently
setx TENANT_SERVICE_PORT "50053"
```

## Common Paths

| Item | Linux/macOS | Windows |
|------|-------------|---------|
| Go binary location | `/usr/local/go/bin/go` | `C:\Program Files\Go\bin\go.exe` |
| GOPATH | `~/go` | `C:\Users\<Username>\go` |
| Go bin (user tools) | `~/go/bin` | `C:\Users\<Username>\go\bin` |
| Project binary | `./bin/tenant-service` | `.\bin\tenant-service.exe` |

## Port Checking

```powershell
# Check if port is in use
netstat -ano | findstr :8083

# Kill process on port (replace PID)
Stop-Process -Id <PID> -Force
```

## Path Separators

- **Windows:** Use backslash `\` in PowerShell or forward slash `/` (both work)
- **Cross-platform Go code:** Use `filepath.Join()` or forward slash `/`

```powershell
# Both work in PowerShell
cd .\internal\service
cd ./internal/service

# In Go code (cross-platform)
path := filepath.Join("internal", "service", "tenant.go")
```

## File Line Endings

```powershell
# Configure Git for Windows (recommended)
git config --global core.autocrlf true

# For this repository only
git config core.autocrlf input
```

## PowerShell Execution Policy

```powershell
# Check current policy
Get-ExecutionPolicy

# Allow scripts (run as Administrator)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run scripts with bypass
powershell -ExecutionPolicy Bypass -File .\build.ps1 build
```

## Docker Commands

```powershell
# Start MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:7.0

# Start Redis
docker run -d -p 6379:6379 --name redis redis:7-alpine

# Stop and remove
docker stop mongodb redis
docker rm mongodb redis

# View logs
docker logs mongodb
docker logs redis
```

## Useful PowerShell Commands

```powershell
# List environment variables
Get-ChildItem Env:

# Find Go version
go version

# Check PowerShell version
$PSVersionTable.PSVersion

# View command history
Get-History

# Clear screen
Clear-Host  # or cls

# Get help for a command
Get-Help <command>
```

## Troubleshooting

### Quick Fixes

```powershell
# Rebuild dependencies
go clean -modcache
.\build.ps1 deps

# Reset Go environment
go env -w GO111MODULE=on
go env -w GOPROXY=https://proxy.golang.org,direct

# Clean build
.\build.ps1 clean
.\build.ps1 build

# Check PATH
$env:PATH -split ';' | Select-String 'go'
```

## Testing the Service

```powershell
# Health check
Invoke-RestMethod http://localhost:8083/health

# Or with curl (via Git Bash)
curl http://localhost:8083/health

# Test with verbose output
Invoke-WebRequest -Uri http://localhost:8083/health -Verbose
```

## VS Code Integration

```powershell
# Open project in VS Code
code .

# Install Go extension
code --install-extension golang.go

# Open specific file
code .\cmd\main.go
```

## Additional Resources

- Full Windows guide: [docs/WINDOWS.md](WINDOWS.md)
- Main README: [README.md](../README.md)
- Dependencies: [docs/DEPENDENCIES.md](DEPENDENCIES.md)
