package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/audio"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/satori/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	AudioHandle = "AudioHandler"
)

type AudioHandler struct {
	// uc represents the usecase interface for authentication.
	uc     audio.AudioUsecase
	logger logger.Logger
	cfg    *config.Config
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAudioHandler(uc audio.AudioUsecase, cfg *config.Config, logr logger.Logger) *AudioHandler {
	return &AudioHandler{uc: uc, cfg: cfg, logger: logr}

}

// audio.AudioUsecase GetTranscription

func (h *AudioHandler) SaveAudio(w http.ResponseWriter, r *http.Request) {

	// ctx.Value("requestId").(string)
	// ctx := r.Context()
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	h.logger.LogInfo(requestId, logger.DeliveryLayer, AudioHandle, "ok")

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, AudioHandle, "ok2")

	file, head, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer file.Close()
	h.logger.LogInfo(requestId, logger.DeliveryLayer, AudioHandle, "ok3")

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	audioDir := h.cfg.AudioUserDir
	filePath := filepath.Join(audioDir, head.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to create file")
		return
	}
	defer out.Close()

	// Копируем содержимое файла из запроса в созданный файл
	_, err = io.Copy(out, file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to save file")
		return
	}

	// Отправляем успешный ответ
	if err := utils.WriteResponse(w, http.StatusOK, map[string]string{
		"message": "File uploaded successfully",
		"path":    filePath,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error write response")
		return
	}

}

func (h *AudioHandler) TranslateAudio(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}
	// Парсим multipart/form-data запрос с ограничением размера файла (5 МБ)
	h.logger.LogDebug("IN Translate Audio")
	h.logger.LogDebug(requestId)

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}

	// Получаем файл из запроса
	file, head, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer file.Close()

	// Проверяем расширение файла
	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	// Создаем временный файл для сохранения аудио
	audioDir := h.cfg.AudioUserDir
	tempFile, err := os.CreateTemp(audioDir, "audio-"+fileType)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to create temp file")
		return
	}
	defer os.Remove(tempFile.Name()) // Удаляем временный файл после завершения

	// Копируем содержимое файла из запроса во временный файл
	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to save temp file")
		return
	}

	// Закрываем временный файл, чтобы убедиться, что данные записаны
	tempFile.Close()

	mlServiceURL := "http://" + h.cfg.MlServer.Address + ":" + h.cfg.MlServer.Port + "/transcribe"
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := sendFileToMLService(client, mlServiceURL, tempFile.Name())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
			utils.WriteError(w, http.StatusGatewayTimeout, "ML service timeout")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "ML service unavailable")
		}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to read ML service response")
		return
	}

	var mlResponse struct {
		Transcription string `json:"transcription"`
	}
	if err := json.Unmarshal(body, &mlResponse); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to parse ML service response")
		return
	}

	// Формируем собственный JSON-ответ
	responseData := map[string]string{
		"transcription": mlResponse.Transcription,
	}

	// Возвращаем транскрипцию клиенту
	if err := utils.WriteResponse(w, http.StatusOK, responseData); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func sendFileToMLService(client *http.Client, url, filePath string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return client.Do(req)
}
