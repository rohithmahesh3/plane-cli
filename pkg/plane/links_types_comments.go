package plane

import "time"

// Link represents an external link attached to an issue
type Link struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedBy string                 `json:"created_by"`
	UpdatedBy string                 `json:"updated_by"`
	Project   string                 `json:"project"`
	Workspace string                 `json:"workspace"`
	Issue     string                 `json:"issue"`
}

// CreateLinkRequest represents a request to create a link
type CreateLinkRequest struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url"`
}

// UpdateLinkRequest represents a request to update a link
type UpdateLinkRequest struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

// IssueType represents a custom issue type
type IssueType struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	LogoProps      map[string]interface{} `json:"logo_props,omitempty"`
	Level          int                    `json:"level"`
	IsActive       bool                   `json:"is_active"`
	IsEpic         bool                   `json:"is_epic"`
	IsDefault      bool                   `json:"is_default"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
	Workspace      string                 `json:"workspace"`
	Project        string                 `json:"project,omitempty"`
	CreatedBy      string                 `json:"created_by"`
	UpdatedBy      string                 `json:"updated_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ExternalID     *string                `json:"external_id,omitempty"`
	ExternalSource *string                `json:"external_source,omitempty"`
}

// CreateIssueTypeRequest represents a request to create an issue type
type CreateIssueTypeRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	LogoProps   map[string]interface{} `json:"logo_props,omitempty"`
	Level       int                    `json:"level,omitempty"`
	IsActive    bool                   `json:"is_active,omitempty"`
	IsEpic      bool                   `json:"is_epic,omitempty"`
	IsDefault   bool                   `json:"is_default,omitempty"`
}

// UpdateIssueTypeRequest represents a request to update an issue type
type UpdateIssueTypeRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	LogoProps   map[string]interface{} `json:"logo_props,omitempty"`
	IsActive    bool                   `json:"is_active,omitempty"`
	IsEpic      bool                   `json:"is_epic,omitempty"`
	IsDefault   bool                   `json:"is_default,omitempty"`
}

// Comment represents a comment on an issue
type Comment struct {
	ID              string    `json:"id"`
	CommentHTML     string    `json:"comment_html"`
	CommentJSON     string    `json:"comment_json,omitempty"`
	CommentStripped string    `json:"comment_stripped,omitempty"`
	Access          string    `json:"access,omitempty"`
	ExternalID      string    `json:"external_id,omitempty"`
	ExternalSource  string    `json:"external_source,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedBy       string    `json:"updated_by"`
	Project         string    `json:"project"`
	Workspace       string    `json:"workspace"`
	Issue           string    `json:"issue"`
	Actor           string    `json:"actor"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	CommentHTML    string `json:"comment_html"`
	CommentJSON    string `json:"comment_json,omitempty"`
	Access         string `json:"access,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	ExternalSource string `json:"external_source,omitempty"`
}

// UpdateCommentRequest represents a request to update a comment
type UpdateCommentRequest struct {
	CommentHTML    string `json:"comment_html,omitempty"`
	CommentJSON    string `json:"comment_json,omitempty"`
	Access         string `json:"access,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	ExternalSource string `json:"external_source,omitempty"`
}
