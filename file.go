package siyuan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	pathpkg "path"
	"strconv"
)

// PutFileRequest writes a file or creates a directory in the workspace.
type PutFileRequest struct {
	Path     string
	IsDir    bool
	ModTime  int64
	FileName string
	File     io.Reader
}

// RenameFileRequest renames a workspace file.
type RenameFileRequest struct {
	Path    string `json:"path"`
	NewPath string `json:"newPath"`
}

// DirEntry describes one workspace directory entry.
type DirEntry struct {
	IsDir     bool   `json:"isDir"`
	IsSymlink bool   `json:"isSymlink"`
	Name      string `json:"name"`
	Updated   int64  `json:"updated"`
}

type filePathRequest struct {
	Path string `json:"path"`
}

// Get returns raw file content from the workspace.
func (s *FileService) Get(ctx context.Context, path string) ([]byte, error) {
	const apiPath = "/api/file/getFile"

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(filePathRequest{Path: path}); err != nil {
		return nil, fmt.Errorf("encode request body for %s: %w", apiPath, err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, apiPath, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("post %s: %w", apiPath, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read file response for %s: %w", apiPath, err)
		}
		return content, nil
	case http.StatusAccepted:
		var envelope responseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			return nil, fmt.Errorf("decode response envelope for %s: %w", apiPath, err)
		}
		return nil, &APIError{
			Code:       envelope.Code,
			Message:    envelope.Msg,
			Path:       apiPath,
			StatusCode: resp.StatusCode,
		}
	default:
		return nil, newHTTPError(apiPath, resp)
	}
}

// Put writes a workspace file or creates a directory.
func (s *FileService) Put(ctx context.Context, req PutFileRequest) error {
	fields := map[string]string{
		"path":  req.Path,
		"isDir": strconv.FormatBool(req.IsDir),
	}
	if req.ModTime != 0 {
		fields["modTime"] = strconv.FormatInt(req.ModTime, 10)
	}

	var files []multipartFile
	if !req.IsDir {
		if req.File == nil {
			return fmt.Errorf("file reader is required when IsDir is false")
		}
		fileName := req.FileName
		if fileName == "" {
			fileName = pathpkg.Base(req.Path)
		}
		if fileName == "." || fileName == "/" {
			fileName = "file"
		}
		files = append(files, multipartFile{
			fieldName: "file",
			fileName:  fileName,
			reader:    req.File,
		})
	}

	return s.client.postMultipart(ctx, "/api/file/putFile", fields, files, nil)
}

// Remove removes a workspace file.
func (s *FileService) Remove(ctx context.Context, path string) error {
	return s.client.postJSON(ctx, "/api/file/removeFile", filePathRequest{Path: path}, nil)
}

// Rename renames a workspace file.
func (s *FileService) Rename(ctx context.Context, req RenameFileRequest) error {
	return s.client.postJSON(ctx, "/api/file/renameFile", req, nil)
}

// ReadDir lists workspace directory entries.
func (s *FileService) ReadDir(ctx context.Context, path string) ([]DirEntry, error) {
	var out []DirEntry
	if err := s.client.postJSON(ctx, "/api/file/readDir", filePathRequest{Path: path}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
