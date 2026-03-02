package context

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	includeAll       bool
	includeCycle     bool
	includeEpic      bool
	includeWorkspace bool
	includeIntake    bool
)

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate CLI command reference for AI agents",
	Long: `Output a concise CLI command reference in markdown format.
Use flags to include additional modules beyond the default set.

Default modules: issue, project, module, state, label, type
Optional modules: --cycle, --epic, --workspace, --intake, --all`,
	RunE: runContext,
}

func init() {
	ContextCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all modules")
	ContextCmd.Flags().BoolVar(&includeCycle, "cycle", false, "Include cycle commands")
	ContextCmd.Flags().BoolVar(&includeEpic, "epic", false, "Include epic commands")
	ContextCmd.Flags().BoolVar(&includeWorkspace, "workspace", false, "Include workspace commands")
	ContextCmd.Flags().BoolVar(&includeIntake, "intake", false, "Include intake commands")
}

func runContext(cmd *cobra.Command, args []string) error {
	output := "# Plane CLI Commands\n\n"

	output += getGlobalFlags()
	output += getIssueCommands()
	output += getProjectCommands()
	output += getModuleCommands()
	output += getStateCommands()
	output += getLabelCommands()
	output += getTypeCommands()

	if includeAll || includeCycle {
		output += getCycleCommands()
	}
	if includeAll || includeEpic {
		output += getEpicCommands()
	}
	if includeAll || includeWorkspace {
		output += getWorkspaceCommands()
	}
	if includeAll || includeIntake {
		output += getIntakeCommands()
	}

	fmt.Print(output)
	return nil
}

func getGlobalFlags() string {
	return `## Global Flags
` + "```" + `
--workspace <slug:text>     Workspace slug (overrides config)
--project <id:text>         Project ID (overrides config)
--output <format>           Output format: json | table | yaml
--no-color                  Disable colored output
--config <path:text>        Config file path
` + "```" + `

`
}

func getIssueCommands() string {
	return `## Issue (aliases: i, issues, ticket)
` + "```" + `
plane-cli issue list [--state <name:text>] [--assignee <@username:text>]
                 [--limit <count:int>]

plane-cli issue view <id:seq_id|uuid>

plane-cli issue create [--title <text>] [--description <text>]
                   [--priority <enum:none|low|medium|high|urgent>]
                   [--assignee <@username:text>...] [--label <text>...]

plane-cli issue edit <id:seq_id|uuid> [--title <text>] [--description <text>]
                 [--priority <enum:none|low|medium|high|urgent>]
                 [--state <enum:backlog|todo|in-progress|done>]
                 [--assignee <@username:text>...] [--label <text>...]

plane-cli issue delete <id:seq_id|uuid>
plane-cli issue search <query:text>

# Issue Comments
plane-cli issue comment list <issue-id:seq_id|uuid>
plane-cli issue comment add <issue-id:seq_id|uuid> [--text <markdown:text>]
                        [--access <enum:INTERNAL|EXTERNAL>]
plane-cli issue comment delete <issue-id:seq_id|uuid> <comment-id:uuid>

# Issue Links
plane-cli issue link list <issue-id:seq_id|uuid>
plane-cli issue link add <issue-id:seq_id|uuid> <url:text> [--title <text>]
plane-cli issue link delete <issue-id:seq_id|uuid> <link-id:uuid>

# Issue Time Tracking
plane-cli issue time list <issue-id:seq_id|uuid>
plane-cli issue time log <issue-id:seq_id|uuid> <duration:minutes|1h30m>
                     [--description <text>]
plane-cli issue time total <issue-id:seq_id|uuid>
plane-cli issue time edit <issue-id:seq_id|uuid> <worklog-id:uuid>
                      [--description <text>] [--duration <minutes|1h30m>]
plane-cli issue time delete <issue-id:seq_id|uuid> <worklog-id:uuid>

# Issue Attachments
plane-cli issue attachment list <issue-id:seq_id|uuid>
plane-cli issue attachment upload <issue-id:seq_id|uuid> <file-path:text>
plane-cli issue attachment edit <issue-id:seq_id|uuid> <attachment-id:uuid>
                           [--name <text>] [--archive | --unarchive]
plane-cli issue attachment delete <issue-id:seq_id|uuid> <attachment-id:uuid>

# Issue Activity
plane-cli issue activity list <issue-id:seq_id|uuid>
plane-cli issue activity view <issue-id:seq_id|uuid> <activity-id:uuid>
` + "```" + `

`
}

func getModuleCommands() string {
	return `## Module (aliases: mod)
` + "```" + `
plane-cli module list [--archived]
plane-cli module view <id:uuid>
plane-cli module create [--name <text>] [--description <text>]
                    [--status <enum:backlog|planned|in-progress|paused|completed|cancelled>]
plane-cli module edit <id:uuid> [--name <text>] [--description <text>] [--status <enum:...>]
plane-cli module delete <id:uuid>
plane-cli module archive <id:uuid>
plane-cli module issues <id:uuid>
plane-cli module add-issues <module-id:uuid> <issue-ids:uuid...>
plane-cli module remove-issue <module-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getStateCommands() string {
	return `## State (aliases: states)
` + "```" + `
plane-cli state list
plane-cli state view <id:uuid>
plane-cli state create [--name <text>] [--description <text>]
                   [--color <hex:#RRGGBB>]
                   [--group <enum:backlog|unstarted|started|completed|cancelled>]
plane-cli state edit <id:uuid> [--name <text>] [--description <text>]
                 [--color <hex>] [--group <enum:...>]
plane-cli state delete <id:uuid>
` + "```" + `

`
}

func getLabelCommands() string {
	return `## Label (aliases: labels, tag)
` + "```" + `
plane-cli label list
plane-cli label view <id:uuid>
plane-cli label create [--name <text>] [--description <text>] [--color <hex:#RRGGBB>]
plane-cli label edit <id:uuid> [--name <text>] [--description <text>] [--color <hex>]
plane-cli label delete <id:uuid>
` + "```" + `

`
}

func getIntakeCommands() string {
	return `## Intake (aliases: inbox, requests)
` + "```" + `
plane-cli intake list
plane-cli intake view <id:uuid>
plane-cli intake create [--name <text>] [--priority <enum:low|medium|high|urgent>]
plane-cli intake delete <id:uuid>
` + "```" + `

`
}

func getTypeCommands() string {
	return `## Type (aliases: issue-type)
` + "```" + `
plane-cli type list
plane-cli type create [--name <text>] [--description <text>]
plane-cli type delete <id:uuid>
` + "```" + `

`
}

func getCycleCommands() string {
	return `## Cycle (aliases: sprint)
` + "```" + `
plane-cli cycle list [--archived]
plane-cli cycle view <id:uuid>
plane-cli cycle create [--name <text>] [--description <text>]
                   [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
plane-cli cycle edit <id:uuid> [--name <text>] [--description <text>]
                 [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
plane-cli cycle delete <id:uuid>
plane-cli cycle archive <id:uuid>
plane-cli cycle issues <id:uuid>
plane-cli cycle add-issues <cycle-id:uuid> <issue-ids:uuid...>
plane-cli cycle remove-issue <cycle-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getEpicCommands() string {
	return `## Epic (aliases: epics)
` + "```" + `
plane-cli epic list
plane-cli epic view <id:uuid>
` + "```" + `

`
}

func getProjectCommands() string {
	return `## Project (aliases: proj)
` + "```" + `
plane-cli project list
plane-cli project create [<name:text>]
plane-cli project info [<id:uuid>]
plane-cli project delete <id:uuid>
plane-cli project members [<id:uuid>]
` + "```" + `

`
}

func getWorkspaceCommands() string {
	return `## Workspace (aliases: ws)
` + "```" + `
plane-cli workspace info [<slug:text>]
plane-cli workspace switch [<slug:text>]
plane-cli workspace members [--search <text>] [--exact] [--limit <count:int>]
` + "```" + `

`
}
