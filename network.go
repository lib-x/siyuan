package siyuan

import "context"

// ForwardProxyRequest sends an HTTP request through SiYuan's forward proxy API.
type ForwardProxyRequest struct {
	URL              string              `json:"url"`
	Method           string              `json:"method,omitempty"`
	Timeout          int                 `json:"timeout,omitempty"`
	ContentType      string              `json:"contentType,omitempty"`
	Headers          []map[string]string `json:"headers,omitempty"`
	Payload          any                 `json:"payload,omitempty"`
	PayloadEncoding  string              `json:"payloadEncoding,omitempty"`
	ResponseEncoding string              `json:"responseEncoding,omitempty"`
}

// ForwardProxyResult describes a forward proxy response.
type ForwardProxyResult struct {
	Body         string         `json:"body"`
	BodyEncoding string         `json:"bodyEncoding"`
	ContentType  string         `json:"contentType"`
	Elapsed      int            `json:"elapsed"`
	Headers      map[string]any `json:"headers"`
	Status       int            `json:"status"`
	URL          string         `json:"url"`
}

// ForwardProxy sends a request through SiYuan's forward proxy API.
func (s *NetworkService) ForwardProxy(ctx context.Context, req ForwardProxyRequest) (*ForwardProxyResult, error) {
	var out ForwardProxyResult
	if err := s.client.postJSON(ctx, "/api/network/forwardProxy", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
