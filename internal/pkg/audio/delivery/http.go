package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/TeaStealers-backend-sem4/internal/pkg/audio"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	SignUpMethod    = "SignUp"
	LoginMethod     = "Login"
	LogoutMethod    = "Logout"
	CheckAuthMethod = "CheckAuth"
)

type AudioHandler struct {
	// uc represents the usecase interface for authentication.
	uc audio.AudioUsecase
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAudioHandler(uc audio.AudioUsecase) *AudioHandler {
	return &AudioHandler{uc: uc}
}

// audio.AudioUsecase GetTranscription

func (h *AudioHandler) SaveAudio(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}

	file, head, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer file.Close()

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	audioDir := "/ouzi/audio"
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
	// Парсим multipart/form-data запрос с ограничением размера файла (5 МБ)
	cfg := config.MustLoad()

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
	audioDir := cfg.AudioUserDir
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

	mlServiceURL := "http://" + cfg.MlServer.Address + ":" + cfg.MlServer.Port + "/transcribe"
	response, err := sendFileToMLService(mlServiceURL, tempFile.Name())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "unable to send file to ML service")
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

func sendFileToMLService(url, filePath string) (*http.Response, error) {
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Создаем multipart/form-data запрос
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// Копируем содержимое файла в запрос
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Закрываем writer, чтобы завершить формирование запроса
	writer.Close()

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос
	client := &http.Client{}
	return client.Do(req)
}
