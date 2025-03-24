package delivery

import (
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"github.com/TeaStealers-backend-sem4/internal/pkg/words"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	WordHandle = "WordHandler"
	GetWord    = "GetWord"
)

type WordHandler struct {
	// uc represents the usecase interface for authentication.
	uc     words.WordUsecase
	cfg    *config.Config
	logger logger.Logger
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewWordHandler(uc words.WordUsecase, cfg *config.Config, logr logger.Logger) *WordHandler {
	return &WordHandler{uc: uc, cfg: cfg, logger: logr}

}

// audio.AudioUsecase GetTranscription

func (h *WordHandler) GetWord(w http.ResponseWriter, r *http.Request) {

	// ctx.Value("requestId").(string)
	// ctx := r.Context()
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	vars := mux.Vars(r)
	word := vars["word"]
	if word == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, GetWord, errors.New("empty word"), 400)
		utils.WriteError(w, http.StatusBadRequest, "word parameter is required")
		return
	}

	audioDir := h.cfg.AudioExampleDir
	audioFilePath := filepath.Join(audioDir, word+".wav") // .wav
	audioFile, err := os.Open(audioFilePath)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "audio not found")
		return
	}
	defer audioFile.Close()

	audioContent, err := io.ReadAll(audioFile)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWord", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "failed to read audio file")
		return
	}

	if err := utils.WriteAudioResponse(w, http.StatusOK, word+".wav", audioContent, word); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWord", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetWord")

}
