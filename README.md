# Plane CLI

A powerful command-line interface for [Plane](https://plane.so) - the open-source project management tool.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- 🔐 **Secure Authentication** - API key stored in OS keyring
- 📊 **Multiple Output Formats** - Table, JSON, and YAML
- 🎯 **Interactive Mode** - Prompts for missing required fields
- 🔍 **Advanced Filtering** - Filter issues by state, priority, assignee, labels, and more
- ⚡ **Fast & Lightweight** - Single binary, no dependencies
- 🎨 **Beautiful Output** - Colored and formatted tables
- 🔧 **Shell Completion** - Bash, Zsh, Fish, and PowerShell support

## Installation

### From Source

```bash
go install github.com/rohithmahesh3/plane-cli@latest
```

### Download Binary

Download the latest release from the [releases page](https://github.com/rohithmahesh3/plane-cli/releases).

### Homebrew (macOS/Linux)

```bash
brew tap rohithmahesh3/plane-cli
brew install plane-cli
```

## Quick Start

### 1. Authenticate

```bash
plane auth login
```

You'll need an API key from Plane:
1. Go to your Plane workspace
2. Navigate to Profile Settings → Personal Access Tokens
3. Create a new token
4. Use it when prompted

### 2. List Workspaces

```bash
plane workspace list
```

### 3. List Projects

```bash
plane project list
```

### 4. List Issues

```bash
plane issue list
```

## Configuration

Configuration is stored in `~/.config/plane-cli/config.yaml`:

```yaml
version: "1.0"
default_workspace: my-workspace
default_project: my-project-id
output_format: table
api_host: https://api.plane.so
```

### Environment Variables

- `PLANE_API_KEY` - Your Plane API key
- `PLANE_WORKSPACE` - Default workspace slug
- `PLANE_PROJECT` - Default project ID

## Usage Examples

### Authentication

```bash
# Interactive login
plane auth login

# Login with flags
plane auth login --token YOUR_API_KEY --workspace my-workspace

# Check authentication status
plane auth status

# Logout
plane auth logout
```

### Workspaces

```bash
# List workspaces
plane workspace list

# Show workspace details
plane workspace info my-workspace

# Switch default workspace
plane workspace switch
```

### Projects

```bash
# List projects
plane project list

# Create a new project
plane project create --name "My Project" --identifier PROJ

# View project details
plane project info PROJECT_ID

# Delete a project
plane project delete PROJECT_ID

# List project members
plane project members PROJECT_ID
```

### Issues (Work Items)

```bash
# List issues
plane issue list

# List with filters
plane issue list --state backlog --priority high
plane issue list --assignee @alice --label bug
plane issue list --cycle "Sprint 1"

# View issue details
plane issue view ISSUE_ID

# Create an issue
plane issue create --title "Bug fix" --priority high
plane issue create -t "Feature request" -d "Description" -p medium

# Edit an issue
plane issue edit ISSUE_ID --state done
plane issue edit ISSUE_ID --assignee @bob

# Delete an issue
plane issue delete ISSUE_ID

# Search issues
plane issue search "login bug"
```

### Output Formats

```bash
# JSON output
plane issue list --output json

# YAML output
plane project list -o yaml

# No colors
plane issue list --no-color
```

## Shell Completion

### Bash

```bash
source <(plane completion bash)
# Add to ~/.bashrc for persistence
```

### Zsh

```bash
source <(plane completion zsh)
# Add to ~/.zshrc for persistence
```

### Fish

```bash
plane completion fish | source
# Save for persistence:
plane completion fish > ~/.config/fish/completions/plane.fish
```

## API Support

This CLI uses the [Plane REST API](https://developers.plane.so/api-reference/introduction) v1.

Supported features:
- ✅ Workspaces (list, info, switch)
- ✅ Projects (list, create, delete, members)
- ✅ Issues/Work Items (list, create, edit, delete, search, filter)
- ⚠️ Cycles (coming soon)
- ⚠️ Modules (coming soon)
- ⚠️ Views (coming soon)

## Development

### Prerequisites

- Go 1.21 or higher
- Git
- pre-commit (optional but recommended)

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Install Locally

```bash
make install
```

### Setup Pre-commit Hooks

We use pre-commit hooks to ensure code quality. Install pre-commit and the hooks:

```bash
# Install pre-commit (if not already installed)
pip install pre-commit

# Install the git hooks
make setup-hooks
```

The pre-commit hooks will automatically run on every commit and check:
- Code formatting (`go fmt`)
- Static analysis (`go vet`)
- Linting (`golangci-lint`)
- Tests (`go test`)

You can also run all checks manually:

```bash
make check  # Runs fmt, vet, lint, and test
```

## Project Structure

```
plane-cli/
├── cmd/                # Command definitions
│   ├── auth/
│   ├── config/
│   ├── cycle/
│   ├── issue/
│   ├── module/
│   ├── project/
│   └── workspace/
├── internal/
│   ├── api/           # API client
│   ├── config/        # Configuration management
│   └── output/        # Output formatting
├── pkg/plane/         # Plane API types
├── main.go
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Plane](https://plane.so) - The open-source project management tool
- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts

## Support

- 🐛 [Report bugs](../../issues)
- 💡 [Request features](../../issues)
- 💬 [Discussions](../../discussions)

---

Made with ❤️ for the Plane community
