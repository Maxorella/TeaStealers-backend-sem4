package delivery

import (
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"github.com/TeaStealers-backend-sem4/internal/pkg/words"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"net/http"
)

const (
	WordHandle = "WordHandler"
	GetWord    = "GetWord"
)

type WordHandler struct {
	// uc represents the usecase interface for authentication.
	uc     words.WordUsecase
	logger logger.Logger
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewWordHandler(uc words.WordUsecase, logr logger.Logger) *WordHandler {
	return &WordHandler{uc: uc, logger: logr}

}

// audio.AudioUsecase GetTranscription

func (h *WordHandler) GetWord(w http.ResponseWriter, r *http.Request) {
	// cfg := config.MustLoad()

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
		h.logger.LogErrorResponse(requestId, utils.DeliveryLayer, GetWord, errors.New("empty word"), 400)
		utils.WriteError(w, http.StatusBadRequest, "word parameter is required")
		return
	}

	// Отправляем успешный ответ
	if err := utils.WriteResponse(w, http.StatusOK, word); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error write response")
		return
	}

}
