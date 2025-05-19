package delivery

import (
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/module"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"net/http"
)

type ModuleHandler struct {
	// uc represents the usecase interface for authentication.
	uc     module.ModuleUsecase
	cfg    *config.Config
	logger logger.Logger
	// minClient *utils2.FileStorageClient
}

func NewModuleHandler(uc module.ModuleUsecase, cfg *config.Config, logr logger.Logger) *ModuleHandler {
	return &ModuleHandler{uc: uc, cfg: cfg, logger: logr}

}

func (h *ModuleHandler) CreateModuleWordHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())
	title := models.ModuleCreate{}

	if err := utils2.ReadRequestData(r, &title); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModuleWordHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	gotId, err := h.uc.CreateModuleWord(r.Context(), title.Title)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModuleWordHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusInternalServerError, "error get all topics")
		return
	}

	newModule := models.ModuleCreate{ID: gotId}

	if err := utils2.WriteResponse(w, http.StatusOK, newModule); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModuleWordHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateModuleWordHandler")
	return

}

func (h *ModuleHandler) CreateModulePhraseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())
	title := models.ModuleCreate{}

	if err := utils2.ReadRequestData(r, &title); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModulePhraseHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	gotId, err := h.uc.CreateModulePhrase(r.Context(), title.Title)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModulePhraseHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusInternalServerError, "error get all topics")
		return
	}

	newModule := models.ModuleCreate{ID: gotId}

	if err := utils2.WriteResponse(w, http.StatusOK, newModule); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateModulePhraseHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateModulePhraseHandler")
	return

}
