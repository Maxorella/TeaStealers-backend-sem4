package delivery

import (
	"github.com/TeaStealers-backend-sem4/internal/pkg/minio"
	"github.com/TeaStealers-backend-sem4/internal/pkg/minio/helpers"
	"github.com/TeaStealers-backend-sem4/internal/pkg/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type Handler struct {
	minioService minio.MinClient
}

func NewMinioHandler(minioService minio.MinClient) *Handler {
	return &Handler{
		minioService: minioService,
	}
}

// CreateOne обработчик для создания одного объекта в хранилище MinIO из переданных данных.
func (h *Handler) CreateOne(w http.ResponseWriter, r *http.Request) {
	// Получаем файл из запроса
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}

	file, head, err := r.FormFile("file") //head _
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

	// Читаем содержимое файла в байтовый срез
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "unable to read form")

		return
	}

	// Создаем структуру FileDataType для хранения данных файла
	fileData := helpers.FileDataType{
		FileName: head.Filename, // Имя файла
		Data:     fileBytes,     // Содержимое файла в виде байтового среза
	}

	// Сохраняем файл в MinIO с помощью метода CreateOne
	link, err := h.minioService.CreateOne(fileData)
	if err != nil {
		// Если не удается сохранить файл, возвращаем ошибку с соответствующим статусом и сообщением
		utils.WriteError(w, http.StatusBadRequest, "unable to save file")
		return
	}

	utils.WriteResponse(w, http.StatusOK, map[string]string{
		"message": "File uploaded successfully",
		"path":    link})
	// Возвращаем успешный ответ с URL-адресом сохраненного файла
	return
}

// GetOne обработчик для получения одного объекта из бакета Minio по его идентификатору.
func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор объекта из параметров URL
	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		utils.WriteError(w, http.StatusBadRequest, "objectID parameter is required")
		return
	}
	// Используем сервис MinIO для получения ссылки на объект
	link, err := h.minioService.GetOne(objectID)
	if err != nil {
		// Если произошла ошибка при получении объекта, возвращаем ошибку с соответствующим статусом и сообщением
		utils.WriteError(w, http.StatusInternalServerError, "Enable to get the object")
		return
	}

	// Возвращаем успешный ответ с URL-адресом полученного файла
	if err = utils.WriteResponse(w, http.StatusOK, link); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to write response")
	}
}

// DeleteOne обработчик для удаления одного объекта из бакета Minio по его идентификатору.
func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectID"]
	if objectID == "" {
		utils.WriteError(w, http.StatusBadRequest, "objectID parameter is required")
		return
	}

	if err := h.minioService.DeleteOne(objectID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Cannot delete the object")
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, "File deleted successfully"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete file")
	}

}

// RegisterRoutes - метод регистрации всех роутов в системе
func (h *Handler) RegisterRoutes(routr *mux.Router) {

	// Здесь мы обозначили все эндпоинты системы с соответствующими хендлерами
	minioRoutes := routr.PathPrefix("/files").Subrouter()
	{
		minioRoutes.HandleFunc("/create", h.CreateOne).Methods(http.MethodPost)
		//	minioRoutes.POST("/", h.minioHandler.CreateOne)
		//	minioRoutes.POST("/many", h.minioHandler.CreateMany)

		minioRoutes.HandleFunc("/get/{objectID}", h.GetOne)

		//	minioRoutes.GET("/:objectID", h.minioHandler.GetOne)
		//	minioRoutes.GET("/many", h.minioHandler.GetMany)

		minioRoutes.HandleFunc("/delete/{objectID}", h.DeleteOne)

		//	minioRoutes.DELETE("/:objectID", h.minioHandler.DeleteOne)
		//	minioRoutes.DELETE("/many", h.minioHandler.DeleteMany)
	}

}
