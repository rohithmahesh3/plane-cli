package inject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetModuleFlags() {
	includeAll = false
	includeProject = false
	includeModule = false
	includeState = false
	includeLabel = false
	includeType = false
	includeCycle = false
	includeWorkspace = false
	includeIntake = false
}

func TestGenerateContentIncludesAllFlagOptions(t *testing.T) {
	resetModuleFlags()
	t.Cleanup(resetModuleFlags)

	content := generateContent()

	assert.Contains(t, content, "`--all` - Include all modules")
	assert.Contains(t, content, "`--project` - Include project commands")
	assert.Contains(t, content, "`--module` - Include module commands")
	assert.Contains(t, content, "`--state` - Include state commands")
	assert.Contains(t, content, "`--label` - Include label commands")
	assert.Contains(t, content, "`--type` - Include type commands")
	assert.Contains(t, content, "`--workspace` - Include workspace commands")
	assert.Contains(t, content, "`--cycle` - Include cycle/sprint commands")
	assert.Contains(t, content, "`--intake` - Include intake commands")
}

func TestInjectCommandDescriptionReflectsCurrentBehavior(t *testing.T) {
	assert.Contains(t, InjectCmd.Long, "Issue command reference (plus optional modules via flags)")
	assert.NotContains(t, InjectCmd.Long, "Quick start commands for issue management")
}
