package siyuan

import "context"

// PushMessageRequest pushes a UI notification.
type PushMessageRequest struct {
	Msg     string `json:"msg"`
	Timeout int    `json:"timeout,omitempty"`
}

// PushedMessage describes a pushed notification.
type PushedMessage struct {
	ID string `json:"id"`
}

// PushMsg pushes an informational message.
func (s *NotificationService) PushMsg(ctx context.Context, req PushMessageRequest) (*PushedMessage, error) {
	var out PushedMessage
	if err := s.client.postJSON(ctx, "/api/notification/pushMsg", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PushErrMsg pushes an error message.
func (s *NotificationService) PushErrMsg(ctx context.Context, req PushMessageRequest) (*PushedMessage, error) {
	var out PushedMessage
	if err := s.client.postJSON(ctx, "/api/notification/pushErrMsg", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
