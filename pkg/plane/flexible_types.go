package plane

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleState can unmarshal from either a string (UUID) or an object
type FlexibleState struct {
	ID          string
	Name        string
	Color       string
	Group       string
	Description string
	IsUUID      bool
}

func (fs *FlexibleState) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		fs.ID = v
		fs.IsUUID = true
		return nil
	case map[string]interface{}:
		fs.ID = getString(v, "id")
		fs.Name = getString(v, "name")
		fs.Color = getString(v, "color")
		fs.Group = getString(v, "group")
		fs.Description = getString(v, "description")
		fs.IsUUID = false
		return nil
	default:
		return fmt.Errorf("state must be string or object, got %T", raw)
	}
}

func (fs *FlexibleState) MarshalJSON() ([]byte, error) {
	if fs.IsUUID {
		return json.Marshal(fs.ID)
	}
	return json.Marshal(map[string]interface{}{
		"id":          fs.ID,
		"name":        fs.Name,
		"color":       fs.Color,
		"group":       fs.Group,
		"description": fs.Description,
	})
}

// FlexibleUser can unmarshal from either a string (UUID) or an object
type FlexibleUser struct {
	ID          string
	Email       string
	DisplayName string
	FirstName   string
	LastName    string
	Avatar      string
	AvatarURL   string
	Role        int
	IsUUID      bool
}

func (fu *FlexibleUser) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		fu.ID = v
		fu.IsUUID = true
		return nil
	case map[string]interface{}:
		fu.ID = getString(v, "id")
		fu.Email = getString(v, "email")
		fu.DisplayName = getString(v, "display_name")
		fu.FirstName = getString(v, "first_name")
		fu.LastName = getString(v, "last_name")
		fu.Avatar = getString(v, "avatar")
		fu.AvatarURL = getString(v, "avatar_url")
		if role, ok := v["role"].(float64); ok {
			fu.Role = int(role)
		}
		fu.IsUUID = false
		return nil
	default:
		return fmt.Errorf("user must be string or object, got %T", raw)
	}
}

func (fu *FlexibleUser) MarshalJSON() ([]byte, error) {
	if fu.IsUUID {
		return json.Marshal(fu.ID)
	}
	return json.Marshal(map[string]interface{}{
		"id":           fu.ID,
		"email":        fu.Email,
		"display_name": fu.DisplayName,
		"first_name":   fu.FirstName,
		"last_name":    fu.LastName,
		"avatar":       fu.Avatar,
		"avatar_url":   fu.AvatarURL,
		"role":         fu.Role,
	})
}

func (fu *FlexibleUser) ToUser() User {
	return User{
		ID:          fu.ID,
		Email:       fu.Email,
		DisplayName: fu.DisplayName,
		FirstName:   fu.FirstName,
		LastName:    fu.LastName,
		Avatar:      fu.Avatar,
		AvatarURL:   fu.AvatarURL,
		Role:        fu.Role,
	}
}

// FlexibleLabel can unmarshal from either a string (UUID) or an object
type FlexibleLabel struct {
	ID          string
	Name        string
	Color       string
	Description string
	IsUUID      bool
}

func (fl *FlexibleLabel) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		fl.ID = v
		fl.IsUUID = true
		return nil
	case map[string]interface{}:
		fl.ID = getString(v, "id")
		fl.Name = getString(v, "name")
		fl.Color = getString(v, "color")
		fl.Description = getString(v, "description")
		fl.IsUUID = false
		return nil
	default:
		return fmt.Errorf("label must be string or object, got %T", raw)
	}
}

func (fl *FlexibleLabel) MarshalJSON() ([]byte, error) {
	if fl.IsUUID {
		return json.Marshal(fl.ID)
	}
	return json.Marshal(map[string]interface{}{
		"id":          fl.ID,
		"name":        fl.Name,
		"color":       fl.Color,
		"description": fl.Description,
	})
}

// ToLabel converts FlexibleLabel to Label
func (fl *FlexibleLabel) ToLabel() Label {
	return Label{
		ID:          fl.ID,
		Name:        fl.Name,
		Color:       fl.Color,
		Description: fl.Description,
	}
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type FlexibleInt64 int64

func (fi *FlexibleInt64) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case float64:
		*fi = FlexibleInt64(int64(v))
		return nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*fi = FlexibleInt64(parsed)
		return nil
	case nil:
		*fi = 0
		return nil
	default:
		return fmt.Errorf("integer must be number or string, got %T", raw)
	}
}

func (fi FlexibleInt64) Int64() int64 {
	return int64(fi)
}

// StateOutput represents a state for JSON/YAML output with state_id and state_name
type StateOutput struct {
	ID   string `json:"state_id"`
	Name string `json:"state_name"`
}

// StateOutputFromIssue builds structured state output from an issue, with fallback
// to state_name when state is returned as a UUID.
func StateOutputFromIssue(issue Issue) StateOutput {
	name := issue.State.Name
	if name == "" {
		name = issue.StateName
	}

	return StateOutput{
		ID:   issue.State.ID,
		Name: name,
	}
}

// String implements fmt.Stringer for table output
func (s StateOutput) String() string {
	if s.Name == "" {
		return "-"
	}
	return s.Name
}
