package dto

import (
	"github.com/weprodev/go-pkg/sanitizer"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

// CreateTemplateRequest is the JSON body for POST /api/v1/workspaces/:wid/templates.
type CreateTemplateRequest struct {
	Name        string `json:"name"`
	UniqueKey   string `json:"unique_key"`
	ChannelType string `json:"channel_type"`
	Category    string `json:"category"`
	Subject     string `json:"subject"`
	ContentHTML string `json:"content_html"`
	IsActive    *bool  `json:"is_active"`
	IsDefault   *bool  `json:"is_default"`
}

// ToDomain maps the request to a domain template ready for create.
func (r CreateTemplateRequest) ToDomain(workspaceID string) *domain.Template {
	t := &domain.Template{
		WorkspaceID: workspaceID,
		Name:        r.Name,
		UniqueKey:   r.UniqueKey,
		ChannelType: r.ChannelType,
		Category:    r.Category,
		Subject:     r.Subject,
		ContentHTML: sanitizer.SanitizeHTML(r.ContentHTML),
		IsActive:    true,
		IsDefault:   false,
	}
	if r.IsActive != nil {
		t.IsActive = *r.IsActive
	}
	if r.IsDefault != nil {
		t.IsDefault = *r.IsDefault
	}
	if t.ChannelType == "" {
		t.ChannelType = "email"
	}
	return t
}

// PatchTemplateRequest is the JSON body for PATCH /api/v1/workspaces/:wid/templates/:tid.
type PatchTemplateRequest struct {
	Name        *string `json:"name"`
	UniqueKey   *string `json:"unique_key"`
	ChannelType *string `json:"channel_type"`
	Category    *string `json:"category"`
	Subject     *string `json:"subject"`
	ContentHTML *string `json:"content_html"`
	IsActive    *bool   `json:"is_active"`
	IsDefault   *bool   `json:"is_default"`
}

// ToPatch maps the request to the service patch command (service loads domain, applies, saves via repository).
func (r PatchTemplateRequest) ToPatch() service.TemplatePatch {
	patch := service.TemplatePatch{
		Name:        r.Name,
		UniqueKey:   r.UniqueKey,
		ChannelType: r.ChannelType,
		Category:    r.Category,
		Subject:     r.Subject,
		ContentHTML: r.ContentHTML,
		IsActive:    r.IsActive,
		IsDefault:   r.IsDefault,
	}
	if r.ContentHTML != nil {
		clean := sanitizer.SanitizeHTML(*r.ContentHTML)
		patch.ContentHTML = &clean
	}
	return patch
}
