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

// - `CAPYBARA` - Capybara
// - `CHAT_GPT` - Chat GPT
// - `CLAUDE` - Claude
// - `COHERE` - Cohere
// - `FALCON` - Falcon
// - `GEMINI` - Gemini
// - `GRANITE` - Granite
// - `GPT3` - GPT-3
// - `FIREWORKS` - Fireworks
// - `LLAMA2` - Llama2
// - `LLAMA3` - Llama3
// - `MISTRAL` - Mistral
// - `MPT` - MPT
// - `OPENCHAT` - OpenChat
// - `PALM` - PaLM
// - `SOLAR` - Solar
// - `TITAN` - Titan
// - `WIZARD` - Wizard
// - `YI` - Yi
// - `ZEPHYR` - Zephyr
type MlModelFamily string

const (
	MlModelFamilyCapybara  MlModelFamily = "CAPYBARA"
	MlModelFamilyChatGpt   MlModelFamily = "CHAT_GPT"
	MlModelFamilyClaude    MlModelFamily = "CLAUDE"
	MlModelFamilyCohere    MlModelFamily = "COHERE"
	MlModelFamilyFalcon    MlModelFamily = "FALCON"
	MlModelFamilyGemini    MlModelFamily = "GEMINI"
	MlModelFamilyGranite   MlModelFamily = "GRANITE"
	MlModelFamilyGpt3      MlModelFamily = "GPT3"
	MlModelFamilyFireworks MlModelFamily = "FIREWORKS"
	MlModelFamilyLlama2    MlModelFamily = "LLAMA2"
	MlModelFamilyLlama3    MlModelFamily = "LLAMA3"
	MlModelFamilyMistral   MlModelFamily = "MISTRAL"
	MlModelFamilyMpt       MlModelFamily = "MPT"
	MlModelFamilyOpenchat  MlModelFamily = "OPENCHAT"
	MlModelFamilyPalm      MlModelFamily = "PALM"
	MlModelFamilySolar     MlModelFamily = "SOLAR"
	MlModelFamilyTitan     MlModelFamily = "TITAN"
	MlModelFamilyWizard    MlModelFamily = "WIZARD"
	MlModelFamilyYi        MlModelFamily = "YI"
	MlModelFamilyZephyr    MlModelFamily = "ZEPHYR"
)

func NewMlModelFamilyFromString(s string) (MlModelFamily, error) {
	switch s {
	case "CAPYBARA":
		return MlModelFamilyCapybara, nil
	case "CHAT_GPT":
		return MlModelFamilyChatGpt, nil
	case "CLAUDE":
		return MlModelFamilyClaude, nil
	case "COHERE":
		return MlModelFamilyCohere, nil
	case "FALCON":
		return MlModelFamilyFalcon, nil
	case "GEMINI":
		return MlModelFamilyGemini, nil
	case "GRANITE":
		return MlModelFamilyGranite, nil
	case "GPT3":
		return MlModelFamilyGpt3, nil
	case "FIREWORKS":
		return MlModelFamilyFireworks, nil
	case "LLAMA2":
		return MlModelFamilyLlama2, nil
	case "LLAMA3":
		return MlModelFamilyLlama3, nil
	case "MISTRAL":
		return MlModelFamilyMistral, nil
	case "MPT":
		return MlModelFamilyMpt, nil
	case "OPENCHAT":
		return MlModelFamilyOpenchat, nil
	case "PALM":
		return MlModelFamilyPalm, nil
	case "SOLAR":
		return MlModelFamilySolar, nil
	case "TITAN":
		return MlModelFamilyTitan, nil
	case "WIZARD":
		return MlModelFamilyWizard, nil
	case "YI":
		return MlModelFamilyYi, nil
	case "ZEPHYR":
		return MlModelFamilyZephyr, nil
	}
	var t MlModelFamily
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (m MlModelFamily) Ptr() *MlModelFamily {
	return &m
}

type MlModelFamilyEnumValueLabel struct {
	Label string        `json:"label"`
	Value MlModelFamily `json:"value,omitempty"`

	_rawJSON json.RawMessage
}

func (m *MlModelFamilyEnumValueLabel) UnmarshalJSON(data []byte) error {
	type unmarshaler MlModelFamilyEnumValueLabel
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*m = MlModelFamilyEnumValueLabel(value)
	m._rawJSON = json.RawMessage(data)
	return nil
}

func (m *MlModelFamilyEnumValueLabel) String() string {
	if len(m._rawJSON) > 0 {
		if value, err := core.StringifyJSON(m._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(m); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", m)
}

type PaginatedMlModelReadList struct {
	Count    *int           `json:"count,omitempty"`
	Next     *string        `json:"next,omitempty"`
	Previous *string        `json:"previous,omitempty"`
	Results  []*MlModelRead `json:"results,omitempty"`

	_rawJSON json.RawMessage
}

func (p *PaginatedMlModelReadList) UnmarshalJSON(data []byte) error {
	type unmarshaler PaginatedMlModelReadList
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*p = PaginatedMlModelReadList(value)
	p._rawJSON = json.RawMessage(data)
	return nil
}

func (p *PaginatedMlModelReadList) String() string {
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
