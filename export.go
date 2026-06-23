package siyuan

import "context"

// MarkdownContent is exported Markdown text.
type MarkdownContent struct {
	HPath   string `json:"hPath"`
	Content string `json:"content"`
}

// ExportResourcesRequest exports workspace files or directories into a zip archive.
type ExportResourcesRequest struct {
	Paths []string `json:"paths"`
	Name  string   `json:"name,omitempty"`
}

// ExportResourcesResult describes an exported resources archive.
type ExportResourcesResult struct {
	Path string `json:"path"`
}

// MarkdownContent exports a document block as Markdown text.
func (s *ExportService) MarkdownContent(ctx context.Context, id BlockID) (*MarkdownContent, error) {
	var out MarkdownContent
	if err := s.client.postJSON(ctx, "/api/export/exportMdContent", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Resources exports workspace files or directories into a zip archive.
func (s *ExportService) Resources(ctx context.Context, req ExportResourcesRequest) (*ExportResourcesResult, error) {
	var out ExportResourcesResult
	if err := s.client.postJSON(ctx, "/api/export/exportResources", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
