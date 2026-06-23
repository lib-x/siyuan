package siyuan

import "context"

// Notebook describes a SiYuan notebook.
type Notebook struct {
	ID     NotebookID `json:"id"`
	Name   string     `json:"name"`
	Icon   string     `json:"icon"`
	Sort   int        `json:"sort"`
	Closed bool       `json:"closed"`
}

// NotebookConf describes notebook-level configuration.
type NotebookConf struct {
	Name                  string `json:"name"`
	Closed                bool   `json:"closed"`
	RefCreateSavePath     string `json:"refCreateSavePath"`
	CreateDocNameTemplate string `json:"createDocNameTemplate"`
	DailyNoteSavePath     string `json:"dailyNoteSavePath"`
	DailyNoteTemplatePath string `json:"dailyNoteTemplatePath"`
}

// NotebookConfResult is returned by GetConf.
type NotebookConfResult struct {
	Box  NotebookID   `json:"box"`
	Conf NotebookConf `json:"conf"`
	Name string       `json:"name"`
}

type notebookIDRequest struct {
	Notebook NotebookID `json:"notebook"`
}

type renameNotebookRequest struct {
	Notebook NotebookID `json:"notebook"`
	Name     string     `json:"name"`
}

type createNotebookRequest struct {
	Name string `json:"name"`
}

type createNotebookResponse struct {
	Notebook Notebook `json:"notebook"`
}

type setNotebookConfRequest struct {
	Notebook NotebookID   `json:"notebook"`
	Conf     NotebookConf `json:"conf"`
}

// List returns all notebooks.
func (s *NotebookService) List(ctx context.Context) ([]Notebook, error) {
	var out struct {
		Notebooks []Notebook `json:"notebooks"`
	}
	if err := s.client.postJSON(ctx, "/api/notebook/lsNotebooks", nil, &out); err != nil {
		return nil, err
	}
	return out.Notebooks, nil
}

// Open opens a notebook by ID.
func (s *NotebookService) Open(ctx context.Context, id NotebookID) error {
	return s.client.postJSON(ctx, "/api/notebook/openNotebook", notebookIDRequest{Notebook: id}, nil)
}

// Close closes a notebook by ID.
func (s *NotebookService) Close(ctx context.Context, id NotebookID) error {
	return s.client.postJSON(ctx, "/api/notebook/closeNotebook", notebookIDRequest{Notebook: id}, nil)
}

// Rename renames a notebook.
func (s *NotebookService) Rename(ctx context.Context, id NotebookID, name string) error {
	return s.client.postJSON(ctx, "/api/notebook/renameNotebook", renameNotebookRequest{Notebook: id, Name: name}, nil)
}

// Create creates a notebook and returns the created notebook.
func (s *NotebookService) Create(ctx context.Context, name string) (*Notebook, error) {
	var out createNotebookResponse
	if err := s.client.postJSON(ctx, "/api/notebook/createNotebook", createNotebookRequest{Name: name}, &out); err != nil {
		return nil, err
	}
	return &out.Notebook, nil
}

// Remove removes a notebook by ID.
func (s *NotebookService) Remove(ctx context.Context, id NotebookID) error {
	return s.client.postJSON(ctx, "/api/notebook/removeNotebook", notebookIDRequest{Notebook: id}, nil)
}

// GetConf returns notebook configuration.
func (s *NotebookService) GetConf(ctx context.Context, id NotebookID) (*NotebookConfResult, error) {
	var out NotebookConfResult
	if err := s.client.postJSON(ctx, "/api/notebook/getNotebookConf", notebookIDRequest{Notebook: id}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetConf saves notebook configuration and returns the saved configuration.
func (s *NotebookService) SetConf(ctx context.Context, id NotebookID, conf NotebookConf) (*NotebookConf, error) {
	var out NotebookConf
	if err := s.client.postJSON(ctx, "/api/notebook/setNotebookConf", setNotebookConfRequest{Notebook: id, Conf: conf}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
