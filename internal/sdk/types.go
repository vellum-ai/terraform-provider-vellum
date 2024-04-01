package api

import (
	json "encoding/json"
	fmt "fmt"
	core "terraform-provider-vellum/internal/sdk/core"
	time "time"
)

// - `ACTIVE` - Active
// - `ARCHIVED` - Archived
type EntityStatus string

const (
	EntityStatusActive   EntityStatus = "ACTIVE"
	EntityStatusArchived EntityStatus = "ARCHIVED"
)

func NewEntityStatusFromString(s string) (EntityStatus, error) {
	switch s {
	case "ACTIVE":
		return EntityStatusActive, nil
	case "ARCHIVED":
		return EntityStatusArchived, nil
	}
	var t EntityStatus
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

// - `DEVELOPMENT` - Development
// - `STAGING` - Staging
// - `PRODUCTION` - Production
type EnvironmentEnum string

const (
	EnvironmentEnumDevelopment EnvironmentEnum = "DEVELOPMENT"
	EnvironmentEnumStaging     EnvironmentEnum = "STAGING"
	EnvironmentEnumProduction  EnvironmentEnum = "PRODUCTION"
)

func NewEnvironmentEnumFromString(s string) (EnvironmentEnum, error) {
	switch s {
	case "DEVELOPMENT":
		return EnvironmentEnumDevelopment, nil
	case "STAGING":
		return EnvironmentEnumStaging, nil
	case "PRODUCTION":
		return EnvironmentEnumProduction, nil
	}
	var t EnvironmentEnum
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

type DocumentIndexRead struct {
	Id      string    `json:"id"`
	Created time.Time `json:"created"`
	// A human-readable label for the document index
	Label string `json:"label"`
	// A name that uniquely identifies this index within its workspace
	Name string `json:"name"`
	// The current status of the document index
	//
	// - `ACTIVE` - Active
	// - `ARCHIVED` - Archived
	Status *EntityStatus `json:"status,omitempty"`
	// The environment this document index is used in
	//
	// - `DEVELOPMENT` - Development
	// - `STAGING` - Staging
	// - `PRODUCTION` - Production
	Environment *EnvironmentEnum `json:"environment,omitempty"`
	// Configuration representing how documents should be indexed
	IndexingConfig map[string]interface{} `json:"indexing_config,omitempty"`

	_rawJSON json.RawMessage
}

func (d *DocumentIndexRead) UnmarshalJSON(data []byte) error {
	type unmarshaler DocumentIndexRead
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*d = DocumentIndexRead(value)
	d._rawJSON = json.RawMessage(data)
	return nil
}

func (d *DocumentIndexRead) String() string {
	if len(d._rawJSON) > 0 {
		if value, err := core.StringifyJSON(d._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(d); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", d)
}

type PaginatedDocumentIndexReadList struct {
	Count    *int                 `json:"count,omitempty"`
	Next     *string              `json:"next,omitempty"`
	Previous *string              `json:"previous,omitempty"`
	Results  []*DocumentIndexRead `json:"results,omitempty"`

	_rawJSON json.RawMessage
}

func (p *PaginatedDocumentIndexReadList) UnmarshalJSON(data []byte) error {
	type unmarshaler PaginatedDocumentIndexReadList
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*p = PaginatedDocumentIndexReadList(value)
	p._rawJSON = json.RawMessage(data)
	return nil
}

func (p *PaginatedDocumentIndexReadList) String() string {
	if len(p._rawJSON) > 0 {
		if value, err := core.StringifyJSON(p._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(p); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", p)
}
