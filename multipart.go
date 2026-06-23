package siyuan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type multipartFile struct {
	fieldName string
	fileName  string
	reader    io.Reader
}

func (c *Client) postMultipart(ctx context.Context, path string, fields map[string]string, files []multipartFile, output any) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return fmt.Errorf("write multipart field %s for %s: %w", name, path, err)
		}
	}

	for _, file := range files {
		if file.reader == nil {
			return fmt.Errorf("multipart file %s reader is nil", file.fileName)
		}
		part, err := writer.CreateFormFile(file.fieldName, file.fileName)
		if err != nil {
			return fmt.Errorf("create multipart file %s for %s: %w", file.fileName, path, err)
		}
		if _, err := io.Copy(part, file.reader); err != nil {
			return fmt.Errorf("copy multipart file %s for %s: %w", file.fileName, path, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer for %s: %w", path, err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post multipart %s: %w", path, err)
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
