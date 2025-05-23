package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	resp, err := json.Marshal(errorResponse)
	if err != nil {
		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) error {
	respSuccess := struct {
		StatusCode int         `json:"statusCode"`
		Message    string      `json:"message,omitempty"`
		Payload    interface{} `json:"payload"`
	}{
		StatusCode: statusCode,
		Payload:    response,
	}
	resp, err := json.Marshal(respSuccess)
	if err != nil {
		return err
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)

	return nil
}

func WriteAudioResponse(w http.ResponseWriter, statusCode int, fileName string, fileContent []byte, text string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile("audio", fileName)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(fileContent)
	if err != nil {
		return err
	}

	_ = writer.WriteField("text", text)

	err = writer.Close()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", writer.FormDataContentType())
	w.WriteHeader(statusCode)

	_, err = io.Copy(w, &body)
	return err
}

func ReadRequestData(r *http.Request, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}
	return nil
}

func ReadResponseData(r *http.Response, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}
	return nil
}
