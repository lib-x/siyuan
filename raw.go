package siyuan

import "context"

// RawService provides low-level access to unwrapped SiYuan JSON API data.
type RawService struct {
	client *Client
}

// Post sends a JSON POST request and decodes the response envelope data into output.
func (s *RawService) Post(ctx context.Context, path string, input any, output any) error {
	return s.client.postJSON(ctx, path, input, output)
}
