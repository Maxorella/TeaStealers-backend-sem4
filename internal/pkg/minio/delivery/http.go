package delivery

import (
	"github.com/TeaStealers-backend-sem4/apperrors"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/minio"
	"github.com/TeaStealers-backend-sem4/internal/pkg/minio/helpers"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type Handler struct {
	minioService minio.MinClient
	logger       logger.Logger
	cfg          *config.Config
}

func NewMinioHandler(minioService minio.MinClient, cfg *config.Config, logr logger.Logger) *Handler {
	return &Handler{
		minioService: minioService,
		cfg:          cfg,
		logger:       logr,
	}
}

func (h *Handler) CreateOne(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrMaxFileSize5, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrMaxFileSize5.Error())
		return
	}

	file, head, err := r.FormFile("file") //head _
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrFileMultipartKeyRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrFileMultipartKeyRequired.Error())
		return
	}

	defer file.Close()

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrWavMp3Only, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrWavMp3Only.Error())
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrReadFileForm, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrReadFileForm.Error())

		return
	}

	fileData := helpers.FileDataType{
		FileName: head.Filename,
		Data:     fileBytes,
	}

	link, err := h.minioService.CreateOne(fileData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrFailedSaveFile, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrFailedSaveFile.Error())
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, link); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", apperrors.ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, apperrors.ErrInternalServer.Error())
	}

}

// GetOne получение одного объекта из Minio по его идентификатору.
func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}

	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", apperrors.ErrObjectIDRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrObjectIDRequired.Error())
		return
	}

	link, err := h.minioService.GetOne(objectID)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", apperrors.ErrFailedToGetFIle, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, apperrors.ErrFailedToGetFIle.Error())
		return
	}

	if err = utils.WriteResponse(w, http.StatusOK, link); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", apperrors.ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, apperrors.ErrInternalServer.Error())
	}
}

// DeleteOne удаление одного объекта из бакета Minio по его идентификатору.
func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		// ctx = context.WithValue(r.Context(), "requestId", requestId)
	}
	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", apperrors.ErrObjectIDRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, apperrors.ErrObjectIDRequired.Error())
		return
	}

	if err := h.minioService.DeleteOne(objectID); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", apperrors.ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, apperrors.ErrInternalServer.Error())
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, "file deleted successfully"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", apperrors.ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, apperrors.ErrInternalServer.Error())
	}

}

func (h *Handler) RegisterRoutes(routr *mux.Router) {
	minioRoutes := routr.PathPrefix("/files").Subrouter()
	{
		minioRoutes.HandleFunc("/create", h.CreateOne).Methods(http.MethodPost)
		minioRoutes.HandleFunc("/get/{objectID}", h.GetOne).Methods(http.MethodGet)
		minioRoutes.HandleFunc("/delete/{objectID}", h.DeleteOne).Methods(http.MethodDelete)
	}
}
