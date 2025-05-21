package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func TranscribeMLService(url string, file multipart.File, filename string, timeout time.Duration) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}
