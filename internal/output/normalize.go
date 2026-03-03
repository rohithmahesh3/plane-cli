package output

import (
	"encoding/json"
	"reflect"
	"strings"
)

func normalizeStructuredOutput(data interface{}) (interface{}, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var normalized interface{}
	if err := json.Unmarshal(payload, &normalized); err != nil {
		return nil, err
	}

	return normalizeValue(normalized), nil
}

func normalizeValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, raw := range typed {
			normalized := normalizeValue(raw)
			if strings.HasSuffix(key, "_html") {
				markdownKey := strings.TrimSuffix(key, "_html") + "_markdown"
				if text, ok := normalized.(string); ok {
					out[markdownKey] = convertHTMLValue(text)
				} else {
					out[markdownKey] = normalized
				}
				continue
			}
			out[key] = normalized
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i := range typed {
			out[i] = normalizeValue(typed[i])
		}
		return out
	default:
		return value
	}
}

type tableField struct {
	header   string
	isHTML   bool
	fieldVal reflect.Value
}

func tableFields(val reflect.Value) []tableField {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make([]tableField, 0, val.NumField())
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tableTag := structField.Tag.Get("table")
		if tableTag == "-" {
			continue
		}

		jsonKey := parseTagName(structField.Tag.Get("json"))
		isHTML := strings.HasSuffix(jsonKey, "_html") || strings.HasSuffix(structField.Name, "HTML")

		header := tableTag
		if header == "" {
			header = structField.Name
		}
		if isHTML {
			header = renameHTMLLabel(header, structField.Name)
		}

		fields = append(fields, tableField{
			header:   header,
			isHTML:   isHTML,
			fieldVal: val.Field(i),
		})
	}

	return fields
}

func parseTagName(tag string) string {
	if tag == "" {
		return ""
	}

	name, _, _ := strings.Cut(tag, ",")
	if name == "-" {
		return ""
	}

	return name
}

func renameHTMLLabel(label, fieldName string) string {
	switch {
	case strings.HasSuffix(label, " HTML"):
		return strings.TrimSuffix(label, " HTML") + " MARKDOWN"
	case strings.HasSuffix(label, "_HTML"):
		return strings.TrimSuffix(label, "_HTML") + "_MARKDOWN"
	case strings.HasSuffix(label, "_html"):
		return strings.TrimSuffix(label, "_html") + "_markdown"
	case strings.HasSuffix(fieldName, "HTML") && strings.HasSuffix(label, "HTML"):
		return strings.TrimSuffix(label, "HTML") + "Markdown"
	default:
		return label
	}
}
