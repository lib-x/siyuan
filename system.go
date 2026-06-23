package siyuan

import "context"

// BootProgress describes SiYuan startup progress.
type BootProgress struct {
	Details  string `json:"details"`
	Progress int    `json:"progress"`
}

// BootProgress returns SiYuan startup progress.
func (s *SystemService) BootProgress(ctx context.Context) (*BootProgress, error) {
	var out BootProgress
	if err := s.client.postJSON(ctx, "/api/system/bootProgress", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Version returns the SiYuan version string.
func (s *SystemService) Version(ctx context.Context) (string, error) {
	var out string
	if err := s.client.postJSON(ctx, "/api/system/version", nil, &out); err != nil {
		return "", err
	}
	return out, nil
}

// CurrentTime returns the current SiYuan server time in milliseconds.
func (s *SystemService) CurrentTime(ctx context.Context) (int64, error) {
	var out int64
	if err := s.client.postJSON(ctx, "/api/system/currentTime", nil, &out); err != nil {
		return 0, err
	}
	return out, nil
}
