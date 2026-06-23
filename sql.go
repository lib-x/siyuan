package siyuan

import "context"

// SQLRow is one dynamic SQL query result row.
type SQLRow map[string]any

type sqlQueryRequest struct {
	Stmt string `json:"stmt"`
}

// Query executes a SQL query.
func (s *SQLService) Query(ctx context.Context, stmt string) ([]SQLRow, error) {
	var out []SQLRow
	if err := s.client.postJSON(ctx, "/api/query/sql", sqlQueryRequest{Stmt: stmt}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FlushTransaction commits pending SQLite transactions.
func (s *SQLService) FlushTransaction(ctx context.Context) error {
	return s.client.postJSON(ctx, "/api/sqlite/flushTransaction", nil, nil)
}
