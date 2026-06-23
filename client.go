package siyuan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultEndpoint is the default SiYuan local API endpoint.
const DefaultEndpoint = "http://127.0.0.1:6806"

const defaultUserAgent = "siyuan-go/0.1.0"

// Client is a SiYuan API client.
type Client struct {
	endpoint   string
	httpClient *http.Client
	token      string
	userAgent  string
	headers    http.Header

	// Notebooks provides notebook APIs.
	Notebooks *NotebookService
	// Documents provides document and file tree APIs.
	Documents *DocumentService
	// Assets provides asset upload APIs.
	Assets *AssetService
	// Blocks provides block manipulation APIs.
	Blocks *BlockService
	// Attributes provides block attribute APIs.
	Attributes *AttributeService
	// SQL provides SQL query and transaction APIs.
	SQL *SQLService
	// Templates provides template rendering APIs.
	Templates *TemplateService
	// Files provides workspace file APIs.
	Files *FileService
	// Export provides export APIs.
	Export *ExportService
	// Convert provides conversion APIs.
	Convert *ConvertService
	// Notifications provides notification APIs.
	Notifications *NotificationService
	// Network provides network proxy APIs.
	Network *NetworkService
	// System provides system APIs.
	System *SystemService
	// Raw provides low-level access to unwrapped JSON API data.
	Raw *RawService
}

// New creates a SiYuan API client.
func New(opts ...Option) (*Client, error) {
	cfg := clientOptions{
		endpoint:   DefaultEndpoint,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  defaultUserAgent,
		headers:    make(http.Header),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	endpoint, err := normalizeEndpoint(cfg.endpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{
		endpoint:   endpoint,
		httpClient: cfg.httpClient,
		token:      cfg.token,
		userAgent:  cfg.userAgent,
		headers:    cfg.headers.Clone(),
	}
	client.Notebooks = &NotebookService{client: client}
	client.Documents = &DocumentService{client: client}
	client.Assets = &AssetService{client: client}
	client.Blocks = &BlockService{client: client}
	client.Attributes = &AttributeService{client: client}
	client.SQL = &SQLService{client: client}
	client.Templates = &TemplateService{client: client}
	client.Files = &FileService{client: client}
	client.Export = &ExportService{client: client}
	client.Convert = &ConvertService{client: client}
	client.Notifications = &NotificationService{client: client}
	client.Network = &NetworkService{client: client}
	client.System = &SystemService{client: client}
	client.Raw = &RawService{client: client}

	return client, nil
}

// Endpoint returns the normalized API endpoint.
func (c *Client) Endpoint() string {
	return c.endpoint
}

func (c *Client) postJSON(ctx context.Context, path string, input any, output any) error {
	var body io.Reader
	if input != nil {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(input); err != nil {
			return fmt.Errorf("encode request body for %s: %w", path, err)
		}
		body = &buf
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newHTTPError(path, resp)
	}

	var envelope responseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode response envelope for %s: %w", path, err)
	}
	if envelope.Code != 0 {
		return &APIError{
			Code:       envelope.Code,
			Message:    envelope.Msg,
			Path:       path,
			StatusCode: resp.StatusCode,
		}
	}
	if output == nil || len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(envelope.Data, output); err != nil {
		return fmt.Errorf("decode response data for %s: %w", path, err)
	}
	return nil
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	endpoint, err := c.resolve(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request for %s: %w", path, err)
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Token "+c.token)
	}
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func (c *Client) resolve(path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("api path must start with slash: %q", path)
	}
	return c.endpoint + path, nil
}

func normalizeEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("endpoint is empty")
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse endpoint: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("endpoint scheme must be http or https")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("endpoint host is empty")
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

type responseEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// NotebookService provides notebook APIs.
type NotebookService struct {
	client *Client
}

// DocumentService provides document and file tree APIs.
type DocumentService struct {
	client *Client
}

// AssetService provides asset upload APIs.
type AssetService struct {
	client *Client
}

// BlockService provides block manipulation APIs.
type BlockService struct {
	client *Client
}

// AttributeService provides block attribute APIs.
type AttributeService struct {
	client *Client
}

// SQLService provides SQL query and transaction APIs.
type SQLService struct {
	client *Client
}

// TemplateService provides template rendering APIs.
type TemplateService struct {
	client *Client
}

// FileService provides workspace file APIs.
type FileService struct {
	client *Client
}

// ExportService provides export APIs.
type ExportService struct {
	client *Client
}

// ConvertService provides conversion APIs.
type ConvertService struct {
	client *Client
}

// NotificationService provides notification APIs.
type NotificationService struct {
	client *Client
}

// NetworkService provides network proxy APIs.
type NetworkService struct {
	client *Client
}

// SystemService provides system APIs.
type SystemService struct {
	client *Client
}
