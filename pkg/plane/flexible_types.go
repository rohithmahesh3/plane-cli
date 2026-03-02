package plane

import (
	"encoding/json"
	"fmt"
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
	Username    string
	Email       string
	DisplayName string
	FirstName   string
	LastName    string
	Avatar      string
	AvatarURL   string
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
		fu.Username = getString(v, "username")
		fu.Email = getString(v, "email")
		fu.DisplayName = getString(v, "display_name")
		fu.FirstName = getString(v, "first_name")
		fu.LastName = getString(v, "last_name")
		fu.Avatar = getString(v, "avatar")
		fu.AvatarURL = getString(v, "avatar_url")
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
		"username":     fu.Username,
		"email":        fu.Email,
		"display_name": fu.DisplayName,
		"first_name":   fu.FirstName,
		"last_name":    fu.LastName,
		"avatar":       fu.Avatar,
		"avatar_url":   fu.AvatarURL,
	})
}

// ToUser converts FlexibleUser to User
func (fu *FlexibleUser) ToUser() User {
	return User{
		ID:          fu.ID,
		Username:    fu.Username,
		Email:       fu.Email,
		DisplayName: fu.DisplayName,
		FirstName:   fu.FirstName,
		LastName:    fu.LastName,
		Avatar:      fu.Avatar,
		AvatarURL:   fu.AvatarURL,
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
