package api

import (
	fmt "fmt"
)

type DocumentIndexCreateRequest struct {
	// A human-readable label for the document index
	Label string `json:"label"`
	// A name that uniquely identifies this index within its workspace
	Name string `json:"name"`
	// The current status of the document index
	//
	// * `ACTIVE` - Active
	// * `ARCHIVED` - Archived
	Status *EntityStatus `json:"status,omitempty"`
	// The environment this document index is used in
	//
	// * `DEVELOPMENT` - Development
	// * `STAGING` - Staging
	// * `PRODUCTION` - Production
	Environment *EnvironmentEnum `json:"environment,omitempty"`
	// Configuration representing how documents should be indexed
	IndexingConfig map[string]interface{} `json:"indexing_config,omitempty"`
	// Optionally specify the id of a document index from which you'd like to copy and re-index its documents into this newly created index
	CopyDocumentsFromIndexId *string `json:"copy_documents_from_index_id,omitempty"`
}

type DocumentIndexesListRequest struct {
	// Number of results to return per page.
	Limit *int `json:"-"`
	// The initial index from which to return the results.
	Offset *int `json:"-"`
	// Which field to use when ordering the results.
	Ordering *string `json:"-"`
	// The current status of the document index
	//
	// - `ACTIVE` - Active
	// - `ARCHIVED` - Archived
	Status *DocumentIndexesListRequestStatus `json:"-"`
}

type PatchedDocumentIndexUpdateRequest struct {
	// A human-readable label for the document index
	Label *string `json:"label,omitempty"`
	// The current status of the document index
	//
	// * `ACTIVE` - Active
	// * `ARCHIVED` - Archived
	Status *EntityStatus `json:"status,omitempty"`
	// The environment this document index is used in
	//
	// * `DEVELOPMENT` - Development
	// * `STAGING` - Staging
	// * `PRODUCTION` - Production
	Environment *EnvironmentEnum `json:"environment,omitempty"`
}

type DocumentIndexesListRequestStatus string

const (
	DocumentIndexesListRequestStatusActive   DocumentIndexesListRequestStatus = "ACTIVE"
	DocumentIndexesListRequestStatusArchived DocumentIndexesListRequestStatus = "ARCHIVED"
)

func NewDocumentIndexesListRequestStatusFromString(s string) (DocumentIndexesListRequestStatus, error) {
	switch s {
	case "ACTIVE":
		return DocumentIndexesListRequestStatusActive, nil
	case "ARCHIVED":
		return DocumentIndexesListRequestStatusArchived, nil
	}
	var t DocumentIndexesListRequestStatus
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (d DocumentIndexesListRequestStatus) Ptr() *DocumentIndexesListRequestStatus {
	return &d
}


type DocumentIndexUpdateRequest struct {
	// A human-readable label for the document index
	Label string `json:"label"`
	// The current status of the document index
	//
	// * `ACTIVE` - Active
	// * `ARCHIVED` - Archived
	Status *EntityStatus `json:"status,omitempty"`
	// The environment this document index is used in
	//
	// * `DEVELOPMENT` - Development
	// * `STAGING` - Staging
	// * `PRODUCTION` - Production
	Environment *EnvironmentEnum `json:"environment,omitempty"`
}