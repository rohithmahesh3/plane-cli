package output

import (
	"encoding/json"
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Formatter struct {
	Format  string
	NoColor bool
	Wide    bool
}

func NewFormatter(format string, noColor bool) *Formatter {
	if noColor {
		color.NoColor = true
	}

	return &Formatter{
		Format:  NormalizeFormat(format),
		NoColor: noColor,
	}
}

func (f *Formatter) Print(data interface{}) error {
	format := NormalizeFormat(f.Format)
	if err := ValidateFormat(format); err != nil {
		return err
	}

	switch format {
	case "json":
		return f.printJSON(data)
	case "yaml":
		return f.printYAML(data)
	}

	return nil
}

func (f *Formatter) printJSON(data interface{}) error {
	normalized, err := normalizeStructuredOutput(data)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(normalized)
}

func (f *Formatter) printYAML(data interface{}) error {
	normalized, err := normalizeStructuredOutput(data)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(os.Stdout)
	defer func() {
		_ = encoder.Close()
	}()
	return encoder.Encode(normalized)
}

func Success(msg string) {
	color.Green("✓ %s", msg)
}

func Error(msg string) {
	color.Red("✗ %s", msg)
}

func Warning(msg string) {
	color.Yellow("⚠ %s", msg)
}

func Info(msg string) {
	color.Cyan("ℹ %s", msg)
}
