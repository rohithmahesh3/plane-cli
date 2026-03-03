package output

import (
	"fmt"
	"strings"
)

const DefaultFormat = "yaml"

func NormalizeFormat(format string) string {
	if strings.TrimSpace(format) == "" {
		return DefaultFormat
	}

	return format
}

func ValidateFormat(format string) error {
	format = NormalizeFormat(format)

	switch format {
	case "json", "yaml":
		return nil
	case "table":
		return fmt.Errorf("invalid output format %q: table output has been removed; supported formats are json, yaml", format)
	default:
		return fmt.Errorf("invalid output format %q: supported formats are json, yaml", format)
	}
}
