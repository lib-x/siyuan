package siyuan

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAssetServiceUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/asset/upload" {
			t.Fatalf("path = %s, want /api/asset/upload", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Token secret" {
			t.Fatalf("Authorization = %q, want Token secret", got)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		if got := r.MultipartForm.Value["assetsDirPath"]; len(got) != 1 || got[0] != "/assets/" {
			t.Fatalf("assetsDirPath = %#v, want /assets/", got)
		}
		files := r.MultipartForm.File["file[]"]
		if len(files) != 2 {
			t.Fatalf("len(files) = %d, want 2", len(files))
		}
		if files[0].Filename != "a.txt" || files[1].Filename != "b.txt" {
			t.Fatalf("filenames = %q/%q, want a.txt/b.txt", files[0].Filename, files[1].Filename)
		}

		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": map[string]any{
				"errFiles": []string{},
				"succMap": map[string]string{
					"a.txt": "assets/a-id.txt",
					"b.txt": "assets/b-id.txt",
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL), WithToken("secret"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := client.Assets.Upload(context.Background(), UploadAssetsRequest{
		AssetsDirPath: "/assets/",
		Files: []UploadAssetFile{
			{Name: "a.txt", Reader: strings.NewReader("a")},
			{Name: "b.txt", Reader: strings.NewReader("b")},
		},
	})
	if err != nil {
		t.Fatalf("Assets.Upload() error = %v", err)
	}
	if result.SuccMap["a.txt"] != "assets/a-id.txt" || result.SuccMap["b.txt"] != "assets/b-id.txt" {
		t.Fatalf("SuccMap = %#v, want uploaded asset paths", result.SuccMap)
	}
}

func TestFileServicePut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/file/putFile" {
			t.Fatalf("path = %s, want /api/file/putFile", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		if got := r.MultipartForm.Value["path"]; len(got) != 1 || got[0] != "/data/foo.txt" {
			t.Fatalf("path field = %#v, want /data/foo.txt", got)
		}
		if got := r.MultipartForm.Value["isDir"]; len(got) != 1 || got[0] != "false" {
			t.Fatalf("isDir field = %#v, want false", got)
		}
		if got := r.MultipartForm.Value["modTime"]; len(got) != 1 || got[0] != "123" {
			t.Fatalf("modTime field = %#v, want 123", got)
		}
		files := r.MultipartForm.File["file"]
		if len(files) != 1 || files[0].Filename != "foo.txt" {
			t.Fatalf("file = %#v, want foo.txt", files)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{"code": 0, "msg": "", "data": nil})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = client.Files.Put(context.Background(), PutFileRequest{
		Path:     "/data/foo.txt",
		ModTime:  123,
		FileName: "foo.txt",
		File:     strings.NewReader("content"),
	})
	if err != nil {
		t.Fatalf("Files.Put() error = %v", err)
	}
}

func TestFileServiceGetReturnsContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/file/getFile" {
			t.Fatalf("path = %s, want /api/file/getFile", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, "file content"); err != nil {
			t.Fatalf("WriteString: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	content, err := client.Files.Get(context.Background(), "/data/foo.txt")
	if err != nil {
		t.Fatalf("Files.Get() error = %v", err)
	}
	if string(content) != "file content" {
		t.Fatalf("content = %q, want file content", string(content))
	}
}

func TestFileServiceGetReturnsAPIErrorForAcceptedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusAccepted, map[string]any{
			"code": 404,
			"msg":  "not found",
			"data": nil,
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Files.Get(context.Background(), "/data/missing.txt")
	if err == nil {
		t.Fatal("Files.Get() error = nil, want APIError")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Files.Get() error = %T, want *APIError", err)
	}
	if apiErr.Code != 404 || apiErr.StatusCode != http.StatusAccepted {
		t.Fatalf("APIError = %#v, want code 404 and status 202", apiErr)
	}
}
