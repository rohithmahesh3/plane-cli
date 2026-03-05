package inject

import (
	"fmt"
	"os"
	"strings"
	"time"

	contextcmd "github.com/rohithmahesh3/plane-cli/cmd/context"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

const (
	markerStart = "<!-- PLANE_TASK_MANAGEMENT_CLI_START -->"
	markerEnd   = "<!-- PLANE_TASK_MANAGEMENT_CLI_END -->"
)

var (
	targetFile string
	dryRun     bool
	force      bool

	// Optional module flags
	includeAll       bool
	includeProject   bool
	includeModule    bool
	includeState     bool
	includeLabel     bool
	includeType      bool
	includeCycle     bool
	includeWorkspace bool
	includeIntake    bool
)

var InjectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject plane-cli context into agent files",
	Long: `Inject plane-cli context documentation into agent files (AGENTS.md, GEMINI.md, CLAUDE.md, CURSOR.md).

This command updates or creates a section in the specified files with:
- Issue command reference (plus optional modules via flags)
- Reference to plane-cli context for full documentation
- Configuration options

The section is marked with HTML comments so it can be automatically updated.

Examples:
  plane-cli inject                    # Update all default agent files
  plane-cli inject --file AGENTS.md   # Update specific file
  plane-cli inject --dry-run          # Show what would change
  plane-cli inject --force            # Force update even if unchanged`,
	RunE: runInject,
}

func init() {
	InjectCmd.Flags().StringVar(&targetFile, "file", "", "Specific file to update (default: all agent files)")
	InjectCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would change without modifying files")
	InjectCmd.Flags().BoolVar(&force, "force", false, "Force update even if content hasn't changed")

	// Optional module flags
	InjectCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all optional modules")
	InjectCmd.Flags().BoolVar(&includeProject, "project", false, "Include project commands")
	InjectCmd.Flags().BoolVar(&includeModule, "module", false, "Include module commands")
	InjectCmd.Flags().BoolVar(&includeState, "state", false, "Include state commands")
	InjectCmd.Flags().BoolVar(&includeLabel, "label", false, "Include label commands")
	InjectCmd.Flags().BoolVar(&includeType, "type", false, "Include type commands")
	InjectCmd.Flags().BoolVar(&includeCycle, "cycle", false, "Include cycle commands")
	InjectCmd.Flags().BoolVar(&includeWorkspace, "workspace", false, "Include workspace commands")
	InjectCmd.Flags().BoolVar(&includeIntake, "intake", false, "Include intake commands")
}

func runInject(cmd *cobra.Command, args []string) error {
	files := getTargetFiles()

	var updated []string
	var skipped []string
	var errors []string

	for _, file := range files {
		result, err := processFile(file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", file, err))
			continue
		}

		switch result {
		case "updated":
			updated = append(updated, file)
		case "skipped":
			skipped = append(skipped, file)
		}
	}

	// Print summary
	fmt.Println()
	if len(updated) > 0 {
		output.Success(fmt.Sprintf("Updated: %s", strings.Join(updated, ", ")))
	}
	if len(skipped) > 0 {
		output.Info(fmt.Sprintf("Skipped (unchanged): %s", strings.Join(skipped, ", ")))
	}
	if len(errors) > 0 {
		for _, err := range errors {
			output.Error(err)
		}
		return fmt.Errorf("failed to update some files")
	}

	return nil
}

func getTargetFiles() []string {
	if targetFile != "" {
		return []string{targetFile}
	}

	return []string{
		"AGENTS.md",
		"GEMINI.md",
		"CLAUDE.md",
		"CURSOR.md",
	}
}

func processFile(filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	newSection := generateContent()
	existingContent := string(content)

	// Check if markers exist
	hasStart := strings.Contains(existingContent, markerStart)
	hasEnd := strings.Contains(existingContent, markerEnd)

	var newContent string

	if hasStart && hasEnd {
		// Update existing section
		newContent = replaceSection(existingContent, newSection)
	} else {
		// Append new section
		newContent = appendSection(existingContent, newSection)
	}

	// Check if content actually changed (unless force)
	if !force && existingContent == newContent {
		return "skipped", nil
	}

	if dryRun {
		fmt.Printf("Would update: %s\n", filePath)
		return "skipped", nil
	}

	// Write updated content
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return "updated", nil
}

func replaceSection(content, newSection string) string {
	startIdx := strings.Index(content, markerStart)
	endIdx := strings.Index(content, markerEnd)

	if startIdx == -1 || endIdx == -1 {
		return content
	}

	endIdx += len(markerEnd)

	return content[:startIdx] + newSection + content[endIdx:]
}

func appendSection(content, newSection string) string {
	// Ensure there's a blank line before the section
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if !strings.HasSuffix(content, "\n\n") {
		content += "\n"
	}

	return content + newSection
}

func generateContent() string {
	timestamp := time.Now().Format("2006-01-02")

	content := fmt.Sprintf(`%s
<!-- Generated by plane-cli inject - Do not edit manually -->
<!-- Last updated: %s -->

## Plane CLI Task Management

The Plane CLI provides command-line access to your Plane workspace for issue tracking, project management, and team collaboration.

`, markerStart, timestamp)

	// Always include full issue commands
	content += contextcmd.GetIssueCommands()

	// Optional modules
	if includeAll || includeProject {
		content += contextcmd.GetProjectCommands()
	}
	if includeAll || includeModule {
		content += contextcmd.GetModuleCommands()
	}
	if includeAll || includeState {
		content += contextcmd.GetStateCommands()
	}
	if includeAll || includeLabel {
		content += contextcmd.GetLabelCommands()
	}
	if includeAll || includeType {
		content += contextcmd.GetTypeCommands()
	}
	if includeAll || includeCycle {
		content += contextcmd.GetCycleCommands()
	}
	if includeAll || includeWorkspace {
		content += contextcmd.GetWorkspaceCommands()
	}
	if includeAll || includeIntake {
		content += contextcmd.GetIntakeCommands()
	}

	// Keep the original sections for reference and tips
	content += `### Full Command Reference

For complete command documentation including modules, states, labels, cycles, and advanced features:

` + "```" + `bash
# Default modules (issue, project, module, state, label, type)
plane-cli context

# All modules including optional (cycle, workspace, intake)
plane-cli context --all

# Specific optional modules
plane-cli context --workspace --cycle --intake --project
` + "```" + `

### Available Context Options

- ` + "`" + `--all` + "`" + ` - Include all modules
- ` + "`" + `--project` + "`" + ` - Include project commands
- ` + "`" + `--module` + "`" + ` - Include module commands
- ` + "`" + `--state` + "`" + ` - Include state commands
- ` + "`" + `--label` + "`" + ` - Include label commands
- ` + "`" + `--type` + "`" + ` - Include type commands
- ` + "`" + `--workspace` + "`" + ` - Include workspace commands
- ` + "`" + `--cycle` + "`" + ` - Include cycle/sprint commands
- ` + "`" + `--intake` + "`" + ` - Include intake commands

### Important Notes

- All entity references (assignees, labels, states) require UUIDs
- Use ` + "`" + `plane-cli workspace members` + "`" + ` to get user IDs
- Use ` + "`" + `plane-cli state list` + "`" + ` to get state IDs
- Use ` + "`" + `plane-cli label list` + "`" + ` to get label IDs

### Tips for Multiline Descriptions

For issues with formatted descriptions, use a heredoc to avoid shell escaping issues:

` + "```" + `bash
plane-cli issue create --title "My issue" --description "$(cat <<'EOF'
## Objective

Detailed description here...

## Requirements
- Item 1
- Item 2
EOF
)"
` + "```" + `

`

	content += markerEnd
	return content
}

// InjectIntoFiles is called from init command
func InjectIntoFiles(files []string) error {
	var errors []string

	for _, file := range files {
		// Skip if file doesn't exist
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		_, err := processFile(file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", file, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during injection: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetDefaultAgentFiles returns the list of default agent files
func GetDefaultAgentFiles() []string {
	return []string{
		"AGENTS.md",
		"GEMINI.md",
		"CLAUDE.md",
		"CURSOR.md",
	}
}
