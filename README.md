# Go CLI Template

A production-ready Go CLI application template with comprehensive CI/CD pipeline, Docker support, and development tooling.

## 🚀 Features

- **Modern Go CLI Structure**: Built with [Cobra](https://github.com/spf13/cobra) for powerful CLI applications
- **Multi-Platform Support**: Build for Linux, macOS, and Windows (AMD64/ARM64)
- **Docker Ready**: Multi-platform Docker images with registry support
- **CI/CD Pipeline**: Complete GitHub Actions workflow with testing, linting, and security scanning
- **Development Tools**: Pre-configured linting, formatting, and testing setup
- **Cursor Rules**: Optimized development experience with Cursor IDE
- **Security First**: Built-in security scanning and vulnerability checks

## 📋 Prerequisites

- Go 1.24+
- Docker (for container builds)
- Make (for build automation)
- Git

## 🛠️ Quick Start

### Using `gonew` (Recommended)

This template is designed to be used with the `gonew` tool for creating new CLI projects:

```bash
# Install gonew if you haven't already
go install golang.org/x/tools/cmd/gonew@latest

# Create a new CLI project from this template
gonew github.com/chengzhycn/go-cli-template your-org/your-cli-app

# Navigate to your new project
cd your-cli-app

# Apply the github ci
mv _github .github

# Initialize git and push to your repository
git init
git add .
git commit -m "Initial commit from go-cli-template"
git remote add origin https://github.com/your-org/your-cli-app.git
git push -u origin main
```

### Manual Setup

If you prefer to clone and customize manually:

```bash
# Clone the template
git clone https://github.com/chengzhycn/go-cli-template.git your-cli-app
cd your-cli-app

# Apply the github ci
mv _github .github

# Remove template-specific files
rm -rf .git
git init

# Update go.mod with your module name
go mod edit -module github.com/your-org/your-cli-app

# Update project references
# Edit main.go, cmd/root.go, and other files to match your project
```

## 🏗️ Project Structure

```
go-cli-template/
├── cmd/                    # CLI commands
│   └── root.go            # Root command
├── internal/              # Private application code
├── pkg/                   # Public libraries
├── .github/               # GitHub Actions workflows
│   └── workflows/
│       └── ci.yml         # CI/CD pipeline
├── Dockerfile             # Multi-platform Docker build
├── .dockerignore          # Docker ignore rules
├── Makefile               # Build automation
├── .golangci.yml          # Linting configuration
├── go.mod                 # Go module definition
├── go.sum                 # Module checksums
├── main.go                # Application entry point
├── LICENSE                # MIT License
└── README.md              # This file
```

## 🔧 Development

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build specific platform
make build-linux
make build-mac
make build-windows
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks (fmt, vet, test, lint)
make check

# Run security checks
make security-check
```

### Docker

```bash
# Build Docker image
make docker-build

# Build multi-platform Docker image
make docker-build-multi

# Run Docker container
make docker-run
```

### Development Workflow

```bash
# Run in development mode
make dev

# Run with file watching (requires air)
make dev-watch

# Install development tools
make install-tools
```

## 🚀 CI/CD Pipeline

The project includes a comprehensive GitHub Actions workflow that:

- **Tests**: Runs tests against multiple Go versions
- **Lints**: Comprehensive code quality checks
- **Builds**: Multi-platform binary builds
- **Scans**: Security vulnerability scanning
- **Releases**: Automatic release creation on tags
- **Docker**: Multi-platform container image publishing

### Pipeline Triggers

- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Version tags (e.g., `v1.0.0`)

### Environment Variables

The CI pipeline uses these environment variables:

- `REGISTRY`: Container registry (default: `ghcr.io`)
- `IMAGE_NAME`: Docker image name (default: repository name)

## 🐳 Docker Support

### Multi-Platform Images

The Dockerfile supports building for multiple architectures:

- Linux AMD64
- Linux ARM64

### Registry Configuration

```bash
# Build with custom registry
REGISTRY=docker.io/myusername VERSION=v1.0.0 make docker-build-multi

# Build for GitHub Container Registry
REGISTRY=ghcr.io/myorg make docker-build-multi
```

### Docker Commands

```bash
# Build and push to registry
make docker-build-multi

# Build locally only
make docker-build-multi-local

# Setup multi-platform builder
make docker-setup-buildx
```

## 🔒 Security

The project includes multiple security measures:

- **Static Analysis**: golangci-lint with security rules
- **Vulnerability Scanning**: govulncheck integration
- **Security Scanning**: gosec for security issues
- **Container Scanning**: Trivy vulnerability scanner
- **Non-root Containers**: Docker images run as non-root user

## 📝 Cursor IDE Integration

This template includes optimized Cursor IDE rules for:

- **Go Best Practices**: Comprehensive Go development guidelines
- **Code Quality**: Automated formatting and linting
- **Security**: Built-in security checks
- **Performance**: Optimization recommendations
- **Testing**: Testing best practices

The rules are automatically applied when using Cursor IDE with this project.

## 🎯 Customization

### Adding Commands

1. Create new command files in `cmd/` directory
2. Follow the Cobra pattern from `cmd/root.go`
3. Register commands in the root command

### Configuration

- Update `cmd/root.go` with your application details
- Modify `main.go` if needed
- Update Docker labels in `Dockerfile`
- Customize CI pipeline in `.github/workflows/ci.yml`

### Environment Variables

The application supports these environment variables:

- `VERSION`: Application version (auto-detected from git)
- `REGISTRY`: Docker registry for builds
- `IMAGE_NAME`: Docker image name

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure quality
5. Submit a pull request

## 📚 Resources

- [Cobra Documentation](https://github.com/spf13/cobra)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Docker Multi-Platform Builds](https://docs.docker.com/build/building/multi-platform/)
- [GitHub Actions](https://docs.github.com/en/actions)

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/chengzhycn/go-cli-template/issues) page
2. Review the documentation above
3. Create a new issue with detailed information

---

**Happy CLI Development! 🎉**
