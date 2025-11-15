# Contributing to atlas-migrate-status

Thank you for your interest in contributing! 🎉

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, database version)

### Suggesting Features

Feature requests are welcome! Please open an issue describing:
- The problem you're trying to solve
- Your proposed solution
- Any alternatives you've considered

### Pull Requests

1. **Fork the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/atlas-migrate-status.git
   cd atlas-migrate-status
   ```

2. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clean, readable code
   - Add tests if applicable
   - Update documentation as needed

4. **Test your changes**
   ```bash
   make test
   make build
   ```

5. **Commit with a clear message**
   ```bash
   git commit -m "Add: brief description of your change"
   ```

6. **Push and create a PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (for testing)

### Building

```bash
make build
```

### Running tests

```bash
make test
```

### Running locally

```bash
go run . --url "postgres://user:pass@localhost:5432/dbname"
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep functions small and focused
- Add comments for complex logic

## Questions?

Feel free to open an issue for any questions!
