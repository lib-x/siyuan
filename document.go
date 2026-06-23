package siyuan

import "context"

// CreateDocWithMarkdownRequest creates a document from Markdown.
type CreateDocWithMarkdownRequest struct {
	Notebook NotebookID `json:"notebook"`
	Path     string     `json:"path"`
	Markdown string     `json:"markdown"`
}

// RenameDocRequest renames a document by storage path.
type RenameDocRequest struct {
	Notebook NotebookID `json:"notebook"`
	Path     string     `json:"path"`
	Title    string     `json:"title"`
}

// RemoveDocRequest removes a document by storage path.
type RemoveDocRequest struct {
	Notebook NotebookID `json:"notebook"`
	Path     string     `json:"path"`
}

// MoveDocsRequest moves documents by storage paths.
type MoveDocsRequest struct {
	FromPaths  []string   `json:"fromPaths"`
	ToNotebook NotebookID `json:"toNotebook"`
	ToPath     string     `json:"toPath"`
}

// MoveDocsByIDRequest moves documents by IDs.
type MoveDocsByIDRequest struct {
	FromIDs []DocumentID `json:"fromIDs"`
	ToID    string       `json:"toID"`
}

// GetHPathByPathRequest resolves a human-readable path from a storage path.
type GetHPathByPathRequest struct {
	Notebook NotebookID `json:"notebook"`
	Path     string     `json:"path"`
}

// GetIDsByHPathRequest resolves block IDs from a human-readable path.
type GetIDsByHPathRequest struct {
	Path     string     `json:"path"`
	Notebook NotebookID `json:"notebook"`
}

// StoragePath describes the notebook and storage path for a block.
type StoragePath struct {
	Notebook NotebookID `json:"notebook"`
	Path     string     `json:"path"`
}

type idRequest struct {
	ID string `json:"id"`
}

type renameDocByIDRequest struct {
	ID    DocumentID `json:"id"`
	Title string     `json:"title"`
}

// CreateWithMarkdown creates a document from Markdown and returns its document ID.
func (s *DocumentService) CreateWithMarkdown(ctx context.Context, req CreateDocWithMarkdownRequest) (DocumentID, error) {
	var out DocumentID
	if err := s.client.postJSON(ctx, "/api/filetree/createDocWithMd", req, &out); err != nil {
		return "", err
	}
	return out, nil
}

// Rename renames a document by storage path.
func (s *DocumentService) Rename(ctx context.Context, req RenameDocRequest) error {
	return s.client.postJSON(ctx, "/api/filetree/renameDoc", req, nil)
}

// RenameByID renames a document by ID.
func (s *DocumentService) RenameByID(ctx context.Context, id DocumentID, title string) error {
	return s.client.postJSON(ctx, "/api/filetree/renameDocByID", renameDocByIDRequest{ID: id, Title: title}, nil)
}

// Remove removes a document by storage path.
func (s *DocumentService) Remove(ctx context.Context, req RemoveDocRequest) error {
	return s.client.postJSON(ctx, "/api/filetree/removeDoc", req, nil)
}

// RemoveByID removes a document by ID.
func (s *DocumentService) RemoveByID(ctx context.Context, id DocumentID) error {
	return s.client.postJSON(ctx, "/api/filetree/removeDocByID", idRequest{ID: string(id)}, nil)
}

// Move moves documents by storage paths.
func (s *DocumentService) Move(ctx context.Context, req MoveDocsRequest) error {
	return s.client.postJSON(ctx, "/api/filetree/moveDocs", req, nil)
}

// MoveByID moves documents by IDs.
func (s *DocumentService) MoveByID(ctx context.Context, req MoveDocsByIDRequest) error {
	return s.client.postJSON(ctx, "/api/filetree/moveDocsByID", req, nil)
}

// GetHPathByPath resolves a human-readable path from a storage path.
func (s *DocumentService) GetHPathByPath(ctx context.Context, req GetHPathByPathRequest) (string, error) {
	var out string
	if err := s.client.postJSON(ctx, "/api/filetree/getHPathByPath", req, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetHPathByID resolves a human-readable path from a block ID.
func (s *DocumentService) GetHPathByID(ctx context.Context, id BlockID) (string, error) {
	var out string
	if err := s.client.postJSON(ctx, "/api/filetree/getHPathByID", idRequest{ID: string(id)}, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetPathByID resolves a storage path from a block ID.
func (s *DocumentService) GetPathByID(ctx context.Context, id BlockID) (*StoragePath, error) {
	var out StoragePath
	if err := s.client.postJSON(ctx, "/api/filetree/getPathByID", idRequest{ID: string(id)}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIDsByHPath resolves block IDs from a human-readable path.
func (s *DocumentService) GetIDsByHPath(ctx context.Context, req GetIDsByHPathRequest) ([]BlockID, error) {
	var out []BlockID
	if err := s.client.postJSON(ctx, "/api/filetree/getIDsByHPath", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
