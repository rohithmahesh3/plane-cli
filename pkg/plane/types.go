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
	ID                  string          `json:"id"`
	SequenceID          int             `json:"sequence_id"`
	Name                string          `json:"name"`
	Description         string          `json:"description,omitempty"`
	DescriptionHTML     string          `json:"description_html,omitempty"`
	DescriptionStripped string          `json:"description_stripped,omitempty"`
	State               FlexibleState   `json:"state"`
	Priority            string          `json:"priority"`
	Assignees           []FlexibleUser  `json:"assignees,omitempty"`
	Labels              []FlexibleLabel `json:"labels,omitempty"`
	CycleID             string          `json:"cycle_id,omitempty"`
	ModuleID            string          `json:"module_id,omitempty"`
	Parent              string          `json:"parent,omitempty"`
	StartDate           string          `json:"start_date,omitempty"`
	TargetDate          string          `json:"target_date,omitempty"`
	EstimatePoint       int             `json:"estimate_point,omitempty"`
	Type                string          `json:"type,omitempty"`
	CompletedAt         *time.Time      `json:"completed_at,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	CreatedBy           string          `json:"created_by,omitempty"`
	UpdatedBy           string          `json:"updated_by,omitempty"`
	ProjectID           string          `json:"project,omitempty"`
	WorkspaceID         string          `json:"workspace,omitempty"`
	IsDraft             bool            `json:"is_draft,omitempty"`
	ArchivedAt          string          `json:"archived_at,omitempty"`
	Sequence            int             `json:"sequence,omitempty"`
	SortOrder           float64         `json:"sort_order,omitempty"`
	StateName           string          `json:"state_name,omitempty"`
	StateGroup          string          `json:"state_group,omitempty"`
	PriorityValue       int             `json:"priority_value,omitempty"`
}

type State struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Color         string    `json:"color"`
	WorkspaceSlug string    `json:"workspace_slug,omitempty"`
	Sequence      float64   `json:"sequence,omitempty"`
	Group         string    `json:"group,omitempty"`
	IsDefault     bool      `json:"default,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type CreateStateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color"`
	Group       string `json:"group,omitempty"`
}

type UpdateStateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Group       string `json:"group,omitempty"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Role        int    `json:"role,omitempty"`
}

type Label struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateLabelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type UpdateLabelRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type Cycle struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitempty"`
	Status      string    `json:"status,omitempty"`
	OwnedBy     string    `json:"owned_by,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ArchivedAt  string    `json:"archived_at,omitempty"`
}

type CreateCycleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateCycleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

type Module struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ArchivedAt  string    `json:"archived_at,omitempty"`
}

type CreateModuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type UpdateModuleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CreateIssueRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	DescriptionHTML string   `json:"description_html,omitempty"`
	State           string   `json:"state,omitempty"`
	Priority        string   `json:"priority,omitempty"`
	Assignees       []string `json:"assignees,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	Parent          string   `json:"parent,omitempty"`
	EstimatePoint   int      `json:"estimate_point,omitempty"`
	Type            string   `json:"type,omitempty"`
	Module          string   `json:"module,omitempty"`
	StartDate       string   `json:"start_date,omitempty"`
	TargetDate      string   `json:"target_date,omitempty"`
}

type UpdateIssueRequest struct {
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	DescriptionHTML string   `json:"description_html,omitempty"`
	State           string   `json:"state,omitempty"`
	Priority        string   `json:"priority,omitempty"`
	Assignees       []string `json:"assignees,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	Parent          string   `json:"parent,omitempty"`
	EstimatePoint   int      `json:"estimate_point,omitempty"`
	Type            string   `json:"type,omitempty"`
	Module          string   `json:"module,omitempty"`
	StartDate       string   `json:"start_date,omitempty"`
	TargetDate      string   `json:"target_date,omitempty"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description,omitempty"`
}

type Worklog struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	ProjectID   string    `json:"project_id,omitempty"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	LoggedBy    string    `json:"logged_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateWorklogRequest struct {
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

type UpdateWorklogRequest struct {
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"`
}

type WorklogTotal struct {
	TotalTime int `json:"total_time"`
}

type Attachment struct {
	ID              string                 `json:"id"`
	Attributes      AttachmentAttributes   `json:"attributes,omitempty"`
	Asset           string                 `json:"asset,omitempty"`
	EntityType      string                 `json:"entity_type,omitempty"`
	IsDeleted       bool                   `json:"is_deleted,omitempty"`
	IsArchived      bool                   `json:"is_archived,omitempty"`
	Size            int64                  `json:"size,omitempty"`
	IsUploaded      bool                   `json:"is_uploaded,omitempty"`
	StorageMetadata map[string]interface{} `json:"storage_metadata,omitempty"`
	CreatedBy       string                 `json:"created_by,omitempty"`
	UpdatedBy       string                 `json:"updated_by,omitempty"`
	Workspace       string                 `json:"workspace,omitempty"`
	Project         string                 `json:"project,omitempty"`
	Issue           string                 `json:"issue,omitempty"`
	CreatedAt       time.Time              `json:"created_at,omitempty"`
	UpdatedAt       time.Time              `json:"updated_at,omitempty"`
}

type AttachmentAttributes struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
	Type string `json:"type,omitempty"`
}

type UploadCredentials struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields,omitempty"`
}

type Epic struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	DescriptionHTML     string    `json:"description_html,omitempty"`
	DescriptionStripped string    `json:"description_stripped,omitempty"`
	State               string    `json:"state,omitempty"`
	Priority            string    `json:"priority,omitempty"`
	Assignees           []string  `json:"assignees,omitempty"`
	Labels              []string  `json:"labels,omitempty"`
	StartDate           string    `json:"start_date,omitempty"`
	TargetDate          string    `json:"target_date,omitempty"`
	SequenceID          int       `json:"sequence_id,omitempty"`
	SortOrder           float64   `json:"sort_order,omitempty"`
	IsDraft             bool      `json:"is_draft,omitempty"`
	CompletedAt         string    `json:"completed_at,omitempty"`
	ArchivedAt          string    `json:"archived_at,omitempty"`
	Project             string    `json:"project,omitempty"`
	Workspace           string    `json:"workspace,omitempty"`
	CreatedBy           string    `json:"created_by,omitempty"`
	UpdatedBy           string    `json:"updated_by,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}

type Page struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	DescriptionHTML string    `json:"description_html,omitempty"`
	CreatedBy       string    `json:"created_by,omitempty"`
	UpdatedBy       string    `json:"updated_by,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type CreatePageRequest struct {
	Name            string `json:"name"`
	DescriptionHTML string `json:"description_html,omitempty"`
}

type UpdatePageRequest struct {
	Name            string `json:"name,omitempty"`
	DescriptionHTML string `json:"description_html,omitempty"`
}

type IntakeIssue struct {
	ID          string    `json:"id"`
	Status      int       `json:"status"`
	SnoozedTill string    `json:"snoozed_till,omitempty"`
	Source      string    `json:"source,omitempty"`
	Inbox       string    `json:"inbox,omitempty"`
	Issue       string    `json:"issue,omitempty"`
	DuplicateTo string    `json:"duplicate_to,omitempty"`
	Project     string    `json:"project,omitempty"`
	Workspace   string    `json:"workspace,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateIntakeIssueRequest struct {
	Issue struct {
		Name     string `json:"name"`
		Priority string `json:"priority,omitempty"`
	} `json:"issue"`
}

type UpdateIntakeIssueRequest struct {
	Status      int    `json:"status,omitempty"`
	SnoozedTill string `json:"snoozed_till,omitempty"`
}

type Activity struct {
	ID            string    `json:"id"`
	Verb          string    `json:"verb"`
	Field         string    `json:"field,omitempty"`
	OldValue      string    `json:"old_value,omitempty"`
	NewValue      string    `json:"new_value,omitempty"`
	Comment       string    `json:"comment,omitempty"`
	Attachments   []string  `json:"attachments,omitempty"`
	OldIdentifier string    `json:"old_identifier,omitempty"`
	NewIdentifier string    `json:"new_identifier,omitempty"`
	Epoch         float64   `json:"epoch,omitempty"`
	Project       string    `json:"project,omitempty"`
	Workspace     string    `json:"workspace,omitempty"`
	Issue         string    `json:"issue,omitempty"`
	IssueComment  string    `json:"issue_comment,omitempty"`
	Actor         string    `json:"actor,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type UpdateAttachmentRequest struct {
	Attributes       AttachmentAttributes   `json:"attributes,omitempty"`
	Asset            string                 `json:"asset,omitempty"`
	EntityType       string                 `json:"entity_type,omitempty"`
	EntityIdentifier string                 `json:"entity_identifier,omitempty"`
	IsDeleted        bool                   `json:"is_deleted,omitempty"`
	IsArchived       bool                   `json:"is_archived,omitempty"`
	ExternalID       string                 `json:"external_id,omitempty"`
	ExternalSource   string                 `json:"external_source,omitempty"`
	Size             int64                  `json:"size,omitempty"`
	IsUploaded       bool                   `json:"is_uploaded,omitempty"`
	StorageMetadata  map[string]interface{} `json:"storage_metadata,omitempty"`
}
