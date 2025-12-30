# Contributing

Thank you for your interest in contributing!

## Development Setup

See README.md for setup instructions.

**For Windows users:** See [docs/WINDOWS.md](docs/WINDOWS.md) for Windows-specific setup instructions.

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Ensure all tests pass
6. Submit a pull request

## Code Style

- Follow Go best practices
- Run `make lint` (or `.\build.ps1 lint` on Windows) before committing
- Add comments for complex logic
- Ensure code is cross-platform compatible

## Testing

- Write unit tests for new features
- Ensure coverage remains high
- Run `make test` (or `.\build.ps1 test` on Windows) before committing

## Cross-Platform Compatibility

When contributing, ensure your code works on all platforms:

- **Use `filepath.Join()`** for file paths instead of hardcoded separators
- **Avoid platform-specific imports** unless absolutely necessary
- **Test on multiple platforms** when possible (Linux, macOS, Windows)
- **Use standard library** cross-platform functions
- **Document platform-specific behavior** if unavoidable


