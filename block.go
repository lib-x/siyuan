package siyuan

import "context"

// InsertBlockRequest inserts a block near another block or under a parent block.
type InsertBlockRequest struct {
	DataType   DataType `json:"dataType"`
	Data       string   `json:"data"`
	NextID     BlockID  `json:"nextID"`
	PreviousID BlockID  `json:"previousID"`
	ParentID   BlockID  `json:"parentID"`
}

// ChildBlockRequest inserts a child block at the start or end of a parent block.
type ChildBlockRequest struct {
	DataType DataType `json:"dataType"`
	Data     string   `json:"data"`
	ParentID BlockID  `json:"parentID"`
}

// UpdateBlockRequest updates a block.
type UpdateBlockRequest struct {
	DataType DataType `json:"dataType"`
	Data     string   `json:"data"`
	ID       BlockID  `json:"id"`
}

// MoveBlockRequest moves a block.
type MoveBlockRequest struct {
	ID         BlockID `json:"id"`
	PreviousID BlockID `json:"previousID"`
	ParentID   BlockID `json:"parentID"`
}

// BlockKramdown contains a block's kramdown source.
type BlockKramdown struct {
	ID       BlockID `json:"id"`
	Kramdown string  `json:"kramdown"`
}

// ChildBlock describes one child block.
type ChildBlock struct {
	ID      BlockID `json:"id"`
	Type    string  `json:"type"`
	SubType string  `json:"subType"`
}

// TransferBlockRefRequest transfers block references from one definition block to another.
type TransferBlockRefRequest struct {
	FromID BlockID   `json:"fromID"`
	ToID   BlockID   `json:"toID"`
	RefIDs []BlockID `json:"refIDs,omitempty"`
}

// Insert inserts a block.
func (s *BlockService) Insert(ctx context.Context, req InsertBlockRequest) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/insertBlock", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Prepend inserts a child block at the start of a parent block.
func (s *BlockService) Prepend(ctx context.Context, req ChildBlockRequest) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/prependBlock", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Append inserts a child block at the end of a parent block.
func (s *BlockService) Append(ctx context.Context, req ChildBlockRequest) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/appendBlock", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a block.
func (s *BlockService) Update(ctx context.Context, req UpdateBlockRequest) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/updateBlock", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a block.
func (s *BlockService) Delete(ctx context.Context, id BlockID) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/deleteBlock", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Move moves a block.
func (s *BlockService) Move(ctx context.Context, req MoveBlockRequest) ([]OperationTransaction, error) {
	var out []OperationTransaction
	if err := s.client.postJSON(ctx, "/api/block/moveBlock", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Fold folds a block.
func (s *BlockService) Fold(ctx context.Context, id BlockID) error {
	return s.client.postJSON(ctx, "/api/block/foldBlock", idRequest{ID: string(id)}, nil)
}

// Unfold unfolds a block.
func (s *BlockService) Unfold(ctx context.Context, id BlockID) error {
	return s.client.postJSON(ctx, "/api/block/unfoldBlock", idRequest{ID: string(id)}, nil)
}

// GetKramdown returns a block's kramdown source.
func (s *BlockService) GetKramdown(ctx context.Context, id BlockID) (*BlockKramdown, error) {
	var out BlockKramdown
	if err := s.client.postJSON(ctx, "/api/block/getBlockKramdown", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetChildBlocks returns child blocks for a parent block.
func (s *BlockService) GetChildBlocks(ctx context.Context, id BlockID) ([]ChildBlock, error) {
	var out []ChildBlock
	if err := s.client.postJSON(ctx, "/api/block/getChildBlocks", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TransferRef transfers block references.
func (s *BlockService) TransferRef(ctx context.Context, req TransferBlockRefRequest) error {
	return s.client.postJSON(ctx, "/api/block/transferBlockRef", req, nil)
}
