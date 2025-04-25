package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/stat"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type AudioHandler struct {
	// uc represents the usecase interface for authentication.
	statUc stat.StatUsecase
	logger logger.Logger
	cfg    *config.Config
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAudioHandler(uc stat.StatUsecase, cfg *config.Config, logr logger.Logger) *AudioHandler {
	return &AudioHandler{statUc: uc, cfg: cfg, logger: logr}
}

func (h *AudioHandler) TranslateAudio(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}

	word := r.FormValue("word")
	if word == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New("no word"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
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
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New("wav and mp3 only"))
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "TranslateAudio", "parsed multiform")

	mlServiceURL := h.cfg.MlServer.TranscribeEndpoint
	timeout := h.cfg.MlServer.Timeout

	response, err := utils.TranscribeMLService(mlServiceURL, file, head.Filename, timeout)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New("ML service timeout"))
			utils.WriteError(w, http.StatusGatewayTimeout, "ML service timeout")
			return
		} else {
			h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New("ML service unavailable"))
			utils.WriteError(w, http.StatusInternalServerError, "ML service unavailable")
			return
		}
	}
	mlAns := models.MlAnswer{}
	err = utils.ReadResponseData(response, mlAns)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New("fail to read ml response"))
		utils.WriteError(w, http.StatusInternalServerError, "fail to read ml response")
		return
	}
	if mlAns.MlError != "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", errors.New(mlAns.MlError))
		utils.WriteError(w, http.StatusInternalServerError, mlAns.MlError)
		return
	}

	result, err := h.statUc.UpdateWordStat(r.Context(), word, mlAns.Transcription)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", err)
		utils.WriteError(w, http.StatusInternalServerError, "update stat error")
		return
	}
	fmt.Printf("ml transcription %s", mlAns.Transcription)
	wordRes := models.WordUserStat{}
	wordRes.Word = word
	wordRes.Transcription = mlAns.Transcription
	wordRes.Progress = &result

	if err := utils.WriteResponse(w, http.StatusOK, wordRes); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "TranslateAudio", err)
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
