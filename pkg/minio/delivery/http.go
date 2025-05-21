package delivery

import (
	"errors"
	"io"
	"net/http"

	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	"github.com/TeaStealers-backend-sem4/pkg/minio"
	"github.com/TeaStealers-backend-sem4/pkg/minio/helpers"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/gorilla/mux"
)

var (
	ErrMaxFileSize5             = errors.New("max file size exceeded (5MB)")
	ErrFileMultipartKeyRequired = errors.New("file key 'file' is required in multipart form")
	ErrReadFileForm             = errors.New("failed to read file from form")
	ErrFailedSaveFile           = errors.New("failed to save file")
	ErrObjectIDRequired         = errors.New("object ID is required")
	ErrFailedToGetFile          = errors.New("failed to get file")
	ErrInternalServer           = errors.New("internal server error")
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
	requestId := utils.GetRequestIDFromCtx(r.Context())
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", ErrMaxFileSize5, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrMaxFileSize5.Error())
		return
	}

	file, head, err := r.FormFile("file")
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", ErrFileMultipartKeyRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrFileMultipartKeyRequired.Error())
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", ErrReadFileForm, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrReadFileForm.Error())
		return
	}

	fileData := helpers.FileDataType{
		FileName: head.Filename,
		Data:     fileBytes,
	}

	link, err := h.minioService.CreateOne(fileData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", ErrFailedSaveFile, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrFailedSaveFile.Error())
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, link); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateOne", ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, ErrInternalServer.Error())
	}
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", ErrObjectIDRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrObjectIDRequired.Error())
		return
	}

	link, err := h.minioService.GetOne(objectID)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", ErrFailedToGetFile, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, ErrFailedToGetFile.Error())
		return
	}

	if err = utils.WriteResponse(w, http.StatusOK, link); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetOne", ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, ErrInternalServer.Error())
	}
}

func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", ErrObjectIDRequired, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, ErrObjectIDRequired.Error())
		return
	}

	if err := h.minioService.DeleteOne(objectID); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, "file deleted successfully"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "DeleteOne", ErrInternalServer, http.StatusBadRequest)
		utils.WriteError(w, http.StatusInternalServerError, ErrInternalServer.Error())
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
