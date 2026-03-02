# Plane CLI

A powerful command-line interface for [Plane](https://plane.so) - the open-source project management tool.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- 🔐 **Secure Authentication** - API key stored in OS keyring
- 📊 **Multiple Output Formats** - Table, JSON, and YAML
- 🎯 **Interactive Mode** - Prompts for missing required fields
- 🔍 **Advanced Filtering** - Filter issues by state, priority, assignee, labels, cycles, and modules
- ⚡ **Fast & Lightweight** - Single binary, no dependencies
- 🎨 **Beautiful Output** - Colored and formatted tables
- 🔧 **Shell Completion** - Bash, Zsh, Fish, and PowerShell support
- ⏱️ **Time Tracking** - Log and manage time spent on issues
- 💬 **Comments & Activity** - Full comment thread and activity history support
- 📎 **Attachments** - Upload and manage file attachments
- 🔗 **Issue Links** - Add external links to issues
- 🤖 **AI Context Generation** - Generate CLI reference for AI agents

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
# Show workspace details
plane workspace info

# Switch default workspace
plane workspace switch my-workspace

# List workspace members
plane workspace members
```

### Projects

```bash
# List projects
plane project list

# List all accessible projects
plane project list --all

# Create a new project
plane project create

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
plane issue list --cycle "Sprint 1" --module "Authentication"

# View issue details (supports sequence ID or UUID)
plane issue view 123
plane issue view uuid-here

# Create an issue
plane issue create --title "Bug fix" --priority high
plane issue create -t "Feature request" -d "Description" -p medium -a @bob

# Edit an issue
plane issue edit 123 --state done
plane issue edit 123 --priority urgent --assignee @alice

# Delete an issue
plane issue delete 123

# Search issues across workspace
plane issue search "login bug"
```

### Issue Comments

```bash
# List comments on an issue
plane issue comment list 123

# Add a comment
plane issue comment add 123 --text "Fixed in PR #42"
plane issue comment add 123 --text "Customer feedback" --access EXTERNAL

# Delete a comment
plane issue comment delete 123 COMMENT_ID
```

### Issue Time Tracking

```bash
# List time logs
plane issue time list 123

# Log time (supports multiple formats)
plane issue time log 123 2h30m --description "Fixed the bug"
plane issue time log 123 90 --description "Code review"

# Show total time logged
plane issue time total 123

# Edit a time log
plane issue time edit 123 WORKLOG_ID --duration 3h

# Delete a time log
plane issue time delete 123 WORKLOG_ID
```

### Issue Links

```bash
# List links on an issue
plane issue link list 123

# Add a link
plane issue link add 123 https://github.com/repo/pull/42 --title "Related PR"

# Delete a link
plane issue link delete 123 LINK_ID
```

### Issue Attachments

```bash
# List attachments
plane issue attachment list 123

# Upload a file
plane issue attachment upload 123 ./screenshot.png

# Edit attachment metadata
plane issue attachment edit 123 ATTACHMENT_ID --name "new-name.png"

# Archive/unarchive
plane issue attachment edit 123 ATTACHMENT_ID --archive

# Delete an attachment
plane issue attachment delete 123 ATTACHMENT_ID
```

### Issue Activity History

```bash
# List activity history
plane issue activity list 123

# View specific activity details
plane issue activity view 123 ACTIVITY_ID
```

### Cycles (Sprints)

```bash
# List cycles
plane cycle list

# List including archived
plane cycle list --archived

# View cycle details
plane cycle view CYCLE_ID

# Create a cycle
plane cycle create --name "Sprint 1" --start-date 2024-01-01 --end-date 2024-01-14

# Edit a cycle
plane cycle edit CYCLE_ID --name "Sprint 1 (Revised)"

# Delete a cycle
plane cycle delete CYCLE_ID

# Archive/unarchive
plane cycle archive CYCLE_ID
plane cycle unarchive CYCLE_ID

# List issues in a cycle
plane cycle issues CYCLE_ID

# Add issues to a cycle
plane cycle add-issues CYCLE_ID ISSUE_ID_1 ISSUE_ID_2

# Remove an issue from a cycle
plane cycle remove-issue CYCLE_ID ISSUE_ID
```

### Modules

```bash
# List modules
plane module list

# List including archived
plane module list --archived

# View module details
plane module view MODULE_ID

# Create a module
plane module create --name "Authentication" --description "Auth features" --status in-progress

# Edit a module
plane module edit MODULE_ID --status completed

# Delete a module
plane module delete MODULE_ID

# Archive/unarchive
plane module archive MODULE_ID
plane module unarchive MODULE_ID

# List issues in a module
plane module issues MODULE_ID

# Add issues to a module
plane module add-issues MODULE_ID ISSUE_ID_1 ISSUE_ID_2

# Remove an issue from a module
plane module remove-issue MODULE_ID ISSUE_ID
```

### Pages (Documentation)

```bash
# List project pages
plane page list

# List workspace pages
plane page list --workspace

# View page details
plane page view PAGE_ID

# Create a page
plane page create --name "API Documentation"
plane page create --name "Team Guidelines" --workspace

# Delete a page
plane page delete PAGE_ID
```

### States (Workflow)

```bash
# List states
plane state list

# View state details
plane state view STATE_ID

# Create a state
plane state create --name "In Review" --color "#F59E0B" --group started

# Edit a state
plane state edit STATE_ID --name "Code Review"

# Delete a state
plane state delete STATE_ID
```

### Labels

```bash
# List labels
plane label list

# View label details
plane label view LABEL_ID

# Create a label
plane label create --name "Bug" --color "#EF4444" --description "Something is broken"

# Edit a label
plane label edit LABEL_ID --color "#3B82F6"

# Delete a label
plane label delete LABEL_ID
```

### Epics

```bash
# List epics
plane epic list

# View epic details
plane epic view EPIC_ID
```

### Intake (Inbox)

```bash
# List intake issues
plane intake list

# View intake issue
plane intake view INTAKE_ID

# Create an intake issue
plane intake create --name "Feature Request" --priority high

# Update intake status
plane intake update INTAKE_ID --status 1  # 1=Accepted

# Delete an intake issue
plane intake delete INTAKE_ID
```

### Issue Types

```bash
# List issue types
plane type list

# Create an issue type
plane type create --name "Bug" --description "Bug reports"

# Delete an issue type
plane type delete TYPE_ID
```

### AI Context Generation

Generate CLI command reference for AI agents:

```bash
# Default modules (issue, module, page, state, label, intake, type)
plane context

# Include additional modules
plane context --cycle
plane context --epic
plane context --project
plane context --workspace

# Include all modules
plane context --all
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

### PowerShell

```powershell
plane completion powershell | Out-String | Invoke-Expression
```

## API Support

This CLI uses the [Plane REST API](https://developers.plane.so/api-reference/introduction) v1.

Supported features:
- ✅ Workspaces (info, switch, members)
- ✅ Projects (list, create, info, delete, members)
- ✅ Issues/Work Items (list, create, edit, delete, search, filter)
- ✅ Issue Comments (list, add, delete)
- ✅ Issue Links (list, add, delete)
- ✅ Issue Time Tracking (list, log, edit, delete, total)
- ✅ Issue Attachments (list, upload, edit, delete)
- ✅ Issue Activity History (list, view)
- ✅ Cycles (list, create, edit, delete, archive, issues management)
- ✅ Modules (list, create, edit, delete, archive, issues management)
- ✅ Pages/Documentation (list, create, view, delete)
- ✅ States/Workflow (list, create, edit, delete)
- ✅ Labels (list, create, edit, delete)
- ✅ Epics (list, view)
- ✅ Intake/Inbox (list, create, view, update, delete)
- ✅ Issue Types (list, create, delete)
- ✅ AI Context Generation

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
│   ├── auth/           # Authentication commands
│   ├── config/         # Config commands
│   ├── context/        # AI context generation
│   ├── cycle/          # Cycle/sprint commands
│   ├── epic/           # Epic commands
│   ├── intake/         # Intake/inbox commands
│   ├── issue/          # Issue commands + subcommands
│   │   ├── activity.go # Activity history
│   │   ├── attachment.go # File attachments
│   │   ├── comment.go  # Comments
│   │   ├── link.go     # External links
│   │   └── time.go     # Time tracking
│   ├── label/          # Label commands
│   ├── module/         # Module commands
│   ├── page/           # Page/documentation commands
│   ├── project/        # Project commands
│   ├── state/          # State/workflow commands
│   ├── type/           # Issue type commands
│   └── workspace/      # Workspace commands
├── internal/
│   ├── api/            # API client and endpoints
│   ├── config/         # Configuration management
│   └── output/         # Output formatting (table/json/yaml)
├── pkg/plane/          # Plane API types and models
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
