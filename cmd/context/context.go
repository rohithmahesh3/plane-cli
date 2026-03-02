package context

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	includeAll       bool
	includeCycle     bool
	includeEpic      bool
	includeProject   bool
	includeWorkspace bool
)

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate CLI command reference for AI agents",
	Long: `Output a concise CLI command reference in markdown format.
Use flags to include additional modules beyond the default set.

Default modules: issue, module, state, label, intake, type
Optional modules: --cycle, --epic, --project, --workspace, --all`,
	RunE: runContext,
}

func init() {
	ContextCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all modules")
	ContextCmd.Flags().BoolVar(&includeCycle, "cycle", false, "Include cycle commands")
	ContextCmd.Flags().BoolVar(&includeEpic, "epic", false, "Include epic commands")
	ContextCmd.Flags().BoolVar(&includeProject, "project", false, "Include project commands")
	ContextCmd.Flags().BoolVar(&includeWorkspace, "workspace", false, "Include workspace commands")
}

func runContext(cmd *cobra.Command, args []string) error {
	output := "# Plane CLI Commands\n\n"

	output += getGlobalFlags()
	output += getIssueCommands()
	output += getModuleCommands()
	output += getStateCommands()
	output += getLabelCommands()
	output += getIntakeCommands()
	output += getTypeCommands()

	if includeAll || includeCycle {
		output += getCycleCommands()
	}
	if includeAll || includeEpic {
		output += getEpicCommands()
	}
	if includeAll || includeProject {
		output += getProjectCommands()
	}
	if includeAll || includeWorkspace {
		output += getWorkspaceCommands()
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
plane issue list [--state <name:text>] [--assignee <@username:text>]
                 [--limit <count:int>]

plane issue view <id:seq_id|uuid>

plane issue create [--title <text>] [--description <text>]
                   [--priority <enum:none|low|medium|high|urgent>]
                   [--assignee <@username:text>...] [--label <text>...]

plane issue edit <id:seq_id|uuid> [--title <text>] [--description <text>]
                 [--priority <enum:none|low|medium|high|urgent>]
                 [--state <enum:backlog|todo|in-progress|done>]
                 [--assignee <@username:text>...] [--label <text>...]

plane issue delete <id:seq_id|uuid>
plane issue search <query:text>

# Issue Comments
plane issue comment list <issue-id:seq_id|uuid>
plane issue comment add <issue-id:seq_id|uuid> [--text <markdown:text>]
                        [--access <enum:INTERNAL|EXTERNAL>]
plane issue comment delete <issue-id:seq_id|uuid> <comment-id:uuid>

# Issue Links
plane issue link list <issue-id:seq_id|uuid>
plane issue link add <issue-id:seq_id|uuid> <url:text> [--title <text>]
plane issue link delete <issue-id:seq_id|uuid> <link-id:uuid>

# Issue Time Tracking
plane issue time list <issue-id:seq_id|uuid>
plane issue time log <issue-id:seq_id|uuid> <duration:minutes|1h30m>
                     [--description <text>]
plane issue time total <issue-id:seq_id|uuid>
plane issue time edit <issue-id:seq_id|uuid> <worklog-id:uuid>
                      [--description <text>] [--duration <minutes|1h30m>]
plane issue time delete <issue-id:seq_id|uuid> <worklog-id:uuid>

# Issue Attachments
plane issue attachment list <issue-id:seq_id|uuid>
plane issue attachment upload <issue-id:seq_id|uuid> <file-path:text>
plane issue attachment edit <issue-id:seq_id|uuid> <attachment-id:uuid>
                           [--name <text>] [--archive | --unarchive]
plane issue attachment delete <issue-id:seq_id|uuid> <attachment-id:uuid>

# Issue Activity
plane issue activity list <issue-id:seq_id|uuid>
plane issue activity view <issue-id:seq_id|uuid> <activity-id:uuid>
` + "```" + `

`
}

func getModuleCommands() string {
	return `## Module (aliases: mod)
` + "```" + `
plane module list [--archived]
plane module view <id:uuid>
plane module create [--name <text>] [--description <text>]
                    [--status <enum:backlog|planned|in-progress|paused|completed|cancelled>]
plane module edit <id:uuid> [--name <text>] [--description <text>] [--status <enum:...>]
plane module delete <id:uuid>
plane module archive <id:uuid>
plane module issues <id:uuid>
plane module add-issues <module-id:uuid> <issue-ids:uuid...>
plane module remove-issue <module-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getStateCommands() string {
	return `## State (aliases: states)
` + "```" + `
plane state list
plane state view <id:uuid>
plane state create [--name <text>] [--description <text>]
                   [--color <hex:#RRGGBB>]
                   [--group <enum:backlog|unstarted|started|completed|cancelled>]
plane state edit <id:uuid> [--name <text>] [--description <text>]
                 [--color <hex>] [--group <enum:...>]
plane state delete <id:uuid>
` + "```" + `

`
}

func getLabelCommands() string {
	return `## Label (aliases: labels, tag)
` + "```" + `
plane label list
plane label view <id:uuid>
plane label create [--name <text>] [--description <text>] [--color <hex:#RRGGBB>]
plane label edit <id:uuid> [--name <text>] [--description <text>] [--color <hex>]
plane label delete <id:uuid>
` + "```" + `

`
}

func getIntakeCommands() string {
	return `## Intake (aliases: inbox, requests)
` + "```" + `
plane intake list
plane intake view <id:uuid>
plane intake create [--name <text>] [--priority <enum:low|medium|high|urgent>]
plane intake delete <id:uuid>
` + "```" + `

`
}

func getTypeCommands() string {
	return `## Type (aliases: issue-type)
` + "```" + `
plane type list
plane type create [--name <text>] [--description <text>]
plane type delete <id:uuid>
` + "```" + `

`
}

func getCycleCommands() string {
	return `## Cycle (aliases: sprint)
` + "```" + `
plane cycle list [--archived]
plane cycle view <id:uuid>
plane cycle create [--name <text>] [--description <text>]
                   [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
plane cycle edit <id:uuid> [--name <text>] [--description <text>]
                 [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
plane cycle delete <id:uuid>
plane cycle archive <id:uuid>
plane cycle issues <id:uuid>
plane cycle add-issues <cycle-id:uuid> <issue-ids:uuid...>
plane cycle remove-issue <cycle-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getEpicCommands() string {
	return `## Epic (aliases: epics)
` + "```" + `
plane epic list
plane epic view <id:uuid>
` + "```" + `

`
}

func getProjectCommands() string {
	return `## Project (aliases: proj)
` + "```" + `
plane project list
plane project create [<name:text>]
plane project info [<id:uuid>]
plane project delete <id:uuid>
plane project members [<id:uuid>]
` + "```" + `

`
}

func getWorkspaceCommands() string {
	return `## Workspace (aliases: ws)
` + "```" + `
plane workspace info [<slug:text>]
plane workspace switch [<slug:text>]
plane workspace members [--search <text>] [--exact] [--limit <count:int>]
` + "```" + `

`
}
