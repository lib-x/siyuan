package siyuan

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewConfiguresDefaultServices(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if client.Endpoint() != DefaultEndpoint {
		t.Fatalf("Endpoint() = %q, want %q", client.Endpoint(), DefaultEndpoint)
	}

	if client.Notebooks == nil || client.Documents == nil || client.Blocks == nil || client.Raw == nil {
		t.Fatal("New() did not initialize service clients")
	}
}

func TestNewRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{name: "empty endpoint", opts: []Option{WithEndpoint("")}},
		{name: "relative endpoint", opts: []Option{WithEndpoint("127.0.0.1:6806")}},
		{name: "unsupported scheme", opts: []Option{WithEndpoint("ftp://127.0.0.1")}},
		{name: "nil http client", opts: []Option{WithHTTPClient(nil)}},
		{name: "empty header name", opts: []Option{WithHeader("", "x")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := New(tt.opts...); err == nil {
				t.Fatal("New() error = nil, want non-nil error")
			}
		})
	}
}

func TestRawPostSendsHeadersAndDecodesData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/api/system/version" {
			t.Fatalf("path = %s, want /api/system/version", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Token secret" {
			t.Fatalf("Authorization = %q, want Token secret", got)
		}
		if got := r.Header.Get("User-Agent"); got != "tests" {
			t.Fatalf("User-Agent = %q, want tests", got)
		}
		if got := r.Header.Get("X-Test"); got != "true" {
			t.Fatalf("X-Test = %q, want true", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if payload["hello"] != "world" {
			t.Fatalf("payload = %#v, want hello=world", payload)
		}

		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": "1.2.3",
		})
	}))
	defer server.Close()

	client, err := New(
		WithEndpoint(server.URL+"/"),
		WithToken("secret"),
		WithUserAgent("tests"),
		WithHeader("X-Test", "true"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	var version string
	err = client.Raw.Post(context.Background(), "/api/system/version", map[string]string{"hello": "world"}, &version)
	if err != nil {
		t.Fatalf("Raw.Post() error = %v", err)
	}
	if version != "1.2.3" {
		t.Fatalf("version = %q, want 1.2.3", version)
	}
}

func TestRawPostReturnsAPIErrorForNonZeroCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 403,
			"msg":  "missing token",
			"data": nil,
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	var out string
	err = client.Raw.Post(context.Background(), "/api/system/version", nil, &out)
	if err == nil {
		t.Fatal("Raw.Post() error = nil, want APIError")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Raw.Post() error = %T, want *APIError", err)
	}
	if apiErr.Code != 403 || apiErr.Message != "missing token" || apiErr.Path != "/api/system/version" {
		t.Fatalf("APIError = %#v, want code/message/path populated", apiErr)
	}
}

func TestRawPostReturnsHTTPErrorForUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = client.Raw.Post(context.Background(), "/api/system/version", nil, nil)
	if err == nil {
		t.Fatal("Raw.Post() error = nil, want HTTPError")
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("Raw.Post() error = %T, want *HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, http.StatusBadGateway)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("Encode response: %v", err)
	}
}
