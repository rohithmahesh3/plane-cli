package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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
		Format:  format,
		NoColor: noColor,
	}
}

func (f *Formatter) Print(data interface{}) error {
	switch f.Format {
	case "json":
		return f.printJSON(data)
	case "yaml":
		return f.printYAML(data)
	case "table", "":
		return f.printTable(data)
	default:
		return fmt.Errorf("unknown output format: %s", f.Format)
	}
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

func (f *Formatter) printTable(data interface{}) error {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Slice {
		return f.printSliceTable(val)
	}

	return f.printStructTable(val)
}

func (f *Formatter) printSliceTable(val reflect.Value) error {
	if val.Len() == 0 {
		fmt.Println("No results found.")
		return nil
	}

	elem := val.Index(0)
	fields := tableFields(elem)
	headers := f.getHeaders(fields)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for i := 0; i < val.Len(); i++ {
		row := f.getRow(tableFields(val.Index(i)))
		table.Append(row)
	}

	table.Render()
	return nil
}

func (f *Formatter) printStructTable(val reflect.Value) error {
	fields := tableFields(val)
	headers := f.getHeaders(fields)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	row := f.getRow(fields)
	table.Append(row)
	table.Render()

	return nil
}

func (f *Formatter) getHeaders(fields []tableField) []string {
	var headers []string
	for _, field := range fields {
		headers = append(headers, strings.ToUpper(field.header[:1])+field.header[1:])
	}

	return headers
}

func (f *Formatter) getRow(fields []tableField) []string {
	var row []string

	for _, field := range fields {
		row = append(row, f.formatValue(field.fieldVal, field.isHTML))
	}

	return row
}

func (f *Formatter) formatValue(val reflect.Value, isHTML bool) string {
	switch val.Kind() {
	case reflect.String:
		s := val.String()
		if isHTML {
			s = convertHTMLValue(s)
		}
		if len(s) > 50 && !f.Wide {
			return s[:47] + "..."
		}
		return s
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Bool:
		return fmt.Sprintf("%v", val.Bool())
	case reflect.Struct:
		if t, ok := val.Interface().(time.Time); ok {
			return f.formatTime(t)
		}
		return fmt.Sprintf("%v", val.Interface())
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return "-"
		}
		var parts []string
		for i := 0; i < val.Len() && i < 3; i++ {
			parts = append(parts, f.formatValue(val.Index(i), false))
		}
		if val.Len() > 3 {
			parts = append(parts, fmt.Sprintf("+%d more", val.Len()-3))
		}
		return strings.Join(parts, ", ")
	case reflect.Ptr:
		if val.IsNil() {
			return "-"
		}
		return f.formatValue(val.Elem(), isHTML)
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

func (f *Formatter) formatTime(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Hour {
		return "just now"
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}

	return t.Format("Jan 2, 2006")
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
