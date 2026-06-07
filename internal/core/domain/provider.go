package domain

import (
	"encoding/json"
	"time"
)

// Provider represents a messaging provider available in the platform catalog.
type Provider struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	ChannelType string                `json:"channel_type"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Fields      []ProviderConfigField `json:"fields,omitempty"`
}

// ProviderConfigField represents a configuration parameter required by a provider.
type ProviderConfigField struct {
	ID           string          `json:"id"`
	ProviderID   string          `json:"provider_id"`
	Key          string          `json:"key"`
	Label        string          `json:"label"`
	Description  string          `json:"description"`
	FieldType    string          `json:"field_type"` // text, password, email, url, etc.
	Required     bool            `json:"required"`
	DefaultValue string          `json:"default_value"`
	Options      json.RawMessage `json:"options,omitempty"`
	SortOrder    int             `json:"sort_order"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
