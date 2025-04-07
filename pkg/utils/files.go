package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type FileStorageClient struct {
	baseURL string
	client  *http.Client
}

func NewFileStorageClient(baseURL string) *FileStorageClient {
	return &FileStorageClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *FileStorageClient) UploadFile(file io.Reader, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/files/create", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	data := struct {
		StatusCode int    `json:"statusCode"`
		Payload    string `json:"payload"`
	}{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", err
	}

	return data.Payload, nil
}

func (c *FileStorageClient) GetFileLink(fileUUID string) (string, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/files/get/"+fileUUID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data := struct {
		StatusCode int    `json:"statusCode"`
		Payload    string `json:"payload"`
	}{}
	respBody, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", err
	}

	if data.Payload == "" {
		return "", fmt.Errorf("empty url in response")
	}

	return data.Payload, nil
}
