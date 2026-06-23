package siyuan

import "context"

// SetBlockAttrsRequest sets attributes for a block.
type SetBlockAttrsRequest struct {
	ID    BlockID           `json:"id"`
	Attrs map[string]string `json:"attrs"`
}

// SetBlockAttrs sets attributes for a block.
func (s *AttributeService) SetBlockAttrs(ctx context.Context, req SetBlockAttrsRequest) error {
	return s.client.postJSON(ctx, "/api/attr/setBlockAttrs", req, nil)
}

// GetBlockAttrs returns attributes for a block.
func (s *AttributeService) GetBlockAttrs(ctx context.Context, id BlockID) (map[string]string, error) {
	var out map[string]string
	if err := s.client.postJSON(ctx, "/api/attr/getBlockAttrs", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
