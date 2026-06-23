package siyuan

import "context"

// RenderTemplateRequest renders a template file for a document.
type RenderTemplateRequest struct {
	ID   DocumentID `json:"id"`
	Path string     `json:"path"`
}

// RenderedTemplate is the rendered template result.
type RenderedTemplate struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}

type renderSprigRequest struct {
	Template string `json:"template"`
}

// Render renders a template file.
func (s *TemplateService) Render(ctx context.Context, req RenderTemplateRequest) (*RenderedTemplate, error) {
	var out RenderedTemplate
	if err := s.client.postJSON(ctx, "/api/template/render", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RenderSprig renders Sprig template content.
func (s *TemplateService) RenderSprig(ctx context.Context, template string) (string, error) {
	var out string
	if err := s.client.postJSON(ctx, "/api/template/renderSprig", renderSprigRequest{Template: template}, &out); err != nil {
		return "", err
	}
	return out, nil
}
