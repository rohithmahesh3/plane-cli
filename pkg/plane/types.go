package plane

import "time"

type Workspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	LogoURL     string    `json:"logo_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description,omitempty"`
	CoverImage  string    `json:"cover_image,omitempty"`
	IconProp    IconProp  `json:"icon_prop,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IconProp struct {
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

type Issue struct {
	ID          string    `json:"id"`
	SequenceID  int       `json:"sequence_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	State       string    `json:"state"`
	Priority    string    `json:"priority"`
	Assignees   []User    `json:"assignees,omitempty"`
	Labels      []string  `json:"labels,omitempty"`
	CycleID     string    `json:"cycle_id,omitempty"`
	ModuleID    string    `json:"module_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type State struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type Cycle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Status      string `json:"status,omitempty"`
}

type Module struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CreateIssueRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Assignees   []string `json:"assignees,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type UpdateIssueRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Assignees   []string `json:"assignees,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description,omitempty"`
}
