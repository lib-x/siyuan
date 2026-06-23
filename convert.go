package siyuan

import "context"

// PandocRequest runs pandoc in a workspace temp conversion directory.
type PandocRequest struct {
	Dir  string   `json:"dir"`
	Args []string `json:"args"`
}

// PandocResult describes the workspace path containing converted files.
type PandocResult struct {
	Path string `json:"path"`
}

// Pandoc runs pandoc with the provided arguments.
func (s *ConvertService) Pandoc(ctx context.Context, req PandocRequest) (*PandocResult, error) {
	var out PandocResult
	if err := s.client.postJSON(ctx, "/api/convert/pandoc", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
