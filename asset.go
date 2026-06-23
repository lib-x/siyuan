package siyuan

import (
	"context"
	"io"
)

// UploadAssetFile is one asset file to upload.
type UploadAssetFile struct {
	Name   string
	Reader io.Reader
}

// UploadAssetsRequest uploads files into an assets directory.
type UploadAssetsRequest struct {
	AssetsDirPath string
	Files         []UploadAssetFile
}

// UploadAssetsResult describes uploaded asset results.
type UploadAssetsResult struct {
	ErrFiles []string          `json:"errFiles"`
	SuccMap  map[string]string `json:"succMap"`
}

// Upload uploads asset files.
func (s *AssetService) Upload(ctx context.Context, req UploadAssetsRequest) (*UploadAssetsResult, error) {
	files := make([]multipartFile, 0, len(req.Files))
	for _, file := range req.Files {
		files = append(files, multipartFile{
			fieldName: "file[]",
			fileName:  file.Name,
			reader:    file.Reader,
		})
	}

	var out UploadAssetsResult
	if err := s.client.postMultipart(ctx, "/api/asset/upload", map[string]string{
		"assetsDirPath": req.AssetsDirPath,
	}, files, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
