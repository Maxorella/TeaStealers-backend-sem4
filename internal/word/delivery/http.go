package delivery

import (
	"context"
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/word"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

const (
	WordHandle = "WordHandler"
	GetWord    = "GetWord"
)

type WordHandler struct {
	// uc represents the usecase interface for authentication.
	uc     word.WordUsecase
	cfg    *config.Config
	logger logger.Logger
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewWordHandler(uc word.WordUsecase, cfg *config.Config, logr logger.Logger) *WordHandler {
	return &WordHandler{uc: uc, cfg: cfg, logger: logr}

}

// audio.AudioUsecase GetTranscription

func (h *WordHandler) GetWord(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
	}

	vars := mux.Vars(r)
	reqWord := vars["word"]
	if reqWord == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, GetWord, errors.New("empty word"), 400)
		utils2.WriteError(w, http.StatusBadRequest, "word parameter is required")
		return
	}
	wordU := &models.WordData{Word: reqWord}
	wordU.Sanitize()
	gotWord, err := h.uc.GetWord(r.Context(), wordU)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get word")
		return
	}

	fileStorageClient := utils2.NewFileStorageClient("http://localhost:8080")
	audioLink, err := fileStorageClient.GetFileLink(gotWord.Link)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetWord", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}
	gotWord.Link = audioLink

	if err := utils2.WriteResponse(w, http.StatusOK, gotWord); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWord", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetWord")

}

func (h *WordHandler) CreateWordHandler(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	data := models.CreateWordData{}

	if err := utils2.ReadRequestData(r, &data); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	_, err := h.uc.CreateWord(context.Background(), &data)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusCreated, "Word created"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error occurred")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateWordHandler")

}

func (h *WordHandler) UploadAudioHandler(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	vars := mux.Vars(r)
	wordVar := vars["word"]
	if wordVar == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UploadAudioHandler", errors.New("empty word"), 400)
		utils2.WriteError(w, http.StatusBadRequest, "word parameter is required")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils2.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "UploadAudioHandler", "parsed multipart form")

	file, head, err := r.FormFile("file")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadAudioHandler", err)
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer file.Close()

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils2.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	h.logger.LogInfo(requestId, logger.DeliveryLayer, "UploadAudioHandler", "got file")

	fileStorageClient := utils2.NewFileStorageClient("http://localhost:8080")
	audioLink, err := fileStorageClient.UploadFile(file, head.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadAudioHandler", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	data := models.WordData{
		Word: wordVar,
		Link: audioLink,
	}
	data.Sanitize()

	if err := h.uc.UploadLink(r.Context(), &data); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadAudioHandler", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to update word audio")
		return
	}
}

func (h *WordHandler) GetRandomWord(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
	}

	vars := mux.Vars(r)
	tag := vars["tag"]

	gotWord, err := h.uc.GetRandomWord(r.Context(), tag)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get word")
		return
	}

	fileStorageClient := utils2.NewFileStorageClient("http://localhost:8080")
	audioLink, err := fileStorageClient.GetFileLink(gotWord.Link)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetRandomWord", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}
	gotWord.Link = audioLink

	if err := utils2.WriteResponse(w, http.StatusOK, gotWord); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetRandomWord", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetRandomWord")

}
