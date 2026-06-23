package siyuan

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError reports a non-zero code returned by the SiYuan API envelope.
type APIError struct {
	Code       int
	Message    string
	Path       string
	StatusCode int
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("siyuan api %s returned code %d", e.Path, e.Code)
	}
	return fmt.Sprintf("siyuan api %s returned code %d: %s", e.Path, e.Code, e.Message)
}

// HTTPError reports a non-successful HTTP response.
type HTTPError struct {
	StatusCode int
	Status     string
	Path       string
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("siyuan api %s returned HTTP %s", e.Path, e.Status)
	}
	return fmt.Sprintf("siyuan api %s returned HTTP %s: %s", e.Path, e.Status, e.Body)
}

func newHTTPError(path string, resp *http.Response) error {
	const maxBody = 4096
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Path:       path,
		Body:       strings.TrimSpace(string(body)),
	}
}
