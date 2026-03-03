package output

import (
	"encoding/json"
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
