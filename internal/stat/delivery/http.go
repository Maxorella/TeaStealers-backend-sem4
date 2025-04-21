package delivery

import (
	"github.com/TeaStealers-backend-sem4/internal/stat"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"net/http"
)

type StatHandler struct {
	// uc represents the usecase interface for authentication.
	uc        stat.StatUsecase
	cfg       *config.Config
	logger    logger.Logger
	minClient *utils2.FileStorageClient
}

func NewStatHandler(uc stat.StatUsecase, cfg *config.Config, logr logger.Logger, minCl *utils2.FileStorageClient) *StatHandler {
	return &StatHandler{uc: uc, cfg: cfg, logger: logr, minClient: minCl}

}

// audio.AudioUsecase GetTranscription
//
//	topic := &models.OneTopic{}
//	err := utils2.ReadRequestData(r, topic)
//	if err != nil {
//		h.logger.LogError(requestId, logger.DeliveryLayer, "GetAllTopics", err)
//		utils2.WriteError(w, http.StatusInternalServerError, "failed to read request")
//		return
//	}
func (h *StatHandler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	gotTopics, err := h.uc.GetAllTopics(r.Context())
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetAllTopics", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusInternalServerError, "error get all topics")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusOK, gotTopics); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetAllTopics", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetAllTopics")
	return
}
