package delivery

import (
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/word"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	"github.com/TeaStealers-backend-sem4/pkg/middleware"
	utils "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	WordHandle = "WordHandler"
	GetWord    = "GetWord"
)

type WordHandler struct {
	ucWord    word.WordUsecase
	cfg       *config.Config
	logger    logger.Logger
	minClient *utils.FileStorageClient
}

func NewWordHandler(ucWord word.WordUsecase, cfg *config.Config, logr logger.Logger, minCl *utils.FileStorageClient) *WordHandler {
	return &WordHandler{ucWord: ucWord, cfg: cfg, logger: logr, minClient: minCl}

}

func (h *WordHandler) CreateWordExerciseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "CreateWordExercise", "parsed multipart form")

	exercise := r.FormValue("exercise")
	if exercise == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	moduleIdStr := r.FormValue("module_id")
	if moduleIdStr == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	moduleId, err := strconv.Atoi(moduleIdStr)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("module_id not int"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
	}
	words := r.FormValue("words")
	if words == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	wordsList := utils.ParseStringArray(words)

	transcriptions := r.FormValue("transcriptions")
	if transcriptions == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	transcriptionsList := utils.ParseStringArray(transcriptions)

	translations := r.FormValue("translations")
	if translations == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	translationsList := utils.ParseStringArray(translations)
	gotId := models.IdStruct{}

	switch exercise {
	case "pronounce":
		audioFile, audioHead, err := r.FormFile("audio")
		if err != nil {
			h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
			utils.WriteError(w, http.StatusBadRequest, "bad data request")
			return
		}
		defer audioFile.Close()
		allowedExtensions := []string{".wav", ".mp3"}
		fileType := strings.ToLower(filepath.Ext(audioHead.Filename))
		if !slices.Contains(allowedExtensions, fileType) {
			utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
			return
		}

		audioLink, err := h.minClient.UploadFile(audioFile, audioHead.Filename)
		if err != nil {
			h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
			utils.WriteError(w, http.StatusInternalServerError, "failed to upload file")
			return
		}

		wordData := models.CreateWordData{Exercise: exercise, ModuleId: &moduleId, Word: wordsList[0], Transcription: transcriptionsList[0], Translation: translationsList[0], AudioLink: audioLink}

		id, err := h.ucWord.CreateWordExercise(r.Context(), &wordData)
		if err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
			utils.WriteError(w, http.StatusInternalServerError, "error create word")
			return
		}
		gotId.Id = &id
	case "pronounceFiew":
		fallthrough
	case "guessWord":
		audioFiles := r.MultipartForm.File["audio"]

		if len(audioFiles) != 2 {
			utils.WriteError(w, http.StatusBadRequest, "not 2 audio uploaded")
			return
		}

		var audioLinks []string
		allowedExtensions := []string{".wav", ".mp3"}

		for _, fileHeader := range audioFiles {
			file, err := fileHeader.Open()
			if err != nil {
				h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
				utils.WriteError(w, http.StatusInternalServerError, "error opening file")
				return
			}
			defer file.Close()

			fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !slices.Contains(allowedExtensions, fileExt) {
				utils.WriteError(w, http.StatusBadRequest, "only .wav and .mp3 allowed")
				return
			}

			audioLink, err := h.minClient.UploadFile(file, fileHeader.Filename)
			if err != nil {
				h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
				utils.WriteError(w, http.StatusInternalServerError, "failed to upload file")
				return
			}
			audioLinks = append(audioLinks, audioLink)
		}

		wordData := models.CreateWordDataList{Exercise: exercise, ModuleId: &moduleId, Word: wordsList,
			Transcription: transcriptionsList, Translation: translationsList, AudioLink: audioLinks}

		id, err := h.ucWord.CreateWordExerciseList(r.Context(), &wordData)
		if err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
			utils.WriteError(w, http.StatusInternalServerError, "error create word")
			return
		}
		gotId.Id = &id
	default:
		utils.WriteError(w, http.StatusBadRequest, "no such exercise")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, gotId); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateWordExercise")
	return
}

func (h *WordHandler) CreatePhraseExerciseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", "parsed multipart form")

	exercise := r.FormValue("exercise")
	if exercise == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	moduleIdStr := r.FormValue("module_id")
	if moduleIdStr == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	moduleId, err := strconv.Atoi(moduleIdStr)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("module_id not int"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
	}

	sentence := r.FormValue("sentence")
	if sentence == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	transcription := r.FormValue("transcription")
	if transcription == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	translate := r.FormValue("translate")
	if translate == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	chain := r.FormValue("chain")
	if chain == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	chainList := utils.ParseStringArray(chain)

	audioFile, audioHead, err := r.FormFile("audio")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audioFile.Close()
	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(audioHead.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	audioLink, err := h.minClient.UploadFile(audioFile, audioHead.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}
	gotId := models.IdStruct{}

	phraseData := models.CreatePhraseData{Exercise: exercise, Sentence: sentence, Transcription: transcription,
		ModuleId:  &moduleId,
		AudioLink: audioLink, Translate: translate, Chain: chainList}

	id, err := h.ucWord.CreatePhraseExercise(r.Context(), &phraseData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}
	gotId.Id = &id

	if err := utils.WriteResponse(w, http.StatusCreated, gotId); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler")
	return
}

func (h *WordHandler) UpdateProgressHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())
	id := r.Context().Value(middleware.CookieName)
	UUID, ok := id.(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}

	progressData := models.ExerciseProgress{UserID: UUID}

	if err := utils.ReadRequestData(r, &progressData); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", err, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	_, err := h.ucWord.CreateUpdateProgress(r.Context(), &progressData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, "Progress saved"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler")
	return
}

func (h *WordHandler) WordModulesHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	gotModules, err := h.ucWord.GetWordModules(r.Context())
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordModulesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, gotModules); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordModulesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "WordModulesHandler")
	return
}

func (h *WordHandler) PhraseModulesHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	gotModules, err := h.ucWord.GetPhraseModules(r.Context())
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordModulesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, gotModules); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordModulesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "WordModulesHandler")
	return
}

func (h *WordHandler) GetWordModuleExercisesHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	vars := mux.Vars(r)
	moduleIDStr := vars["id"]

	moduleID, err := strconv.Atoi(moduleIDStr)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWordModuleExercisesHandler", err, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, "invalid module ID format")
		return
	}

	if moduleID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "module ID must be positive")
		return
	}
	id := r.Context().Value(middleware.CookieName)
	userId, ok := id.(string)
	if !ok {
		userId = ""
	}

	gotModules, err := h.ucWord.GetWordModuleExercises(r.Context(), userId, moduleID)

	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWordModuleExercisesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, gotModules); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWordModuleExercisesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetWordModuleExercisesHandler")
	return
}

func (h *WordHandler) GetPhraseModuleExercisesHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	vars := mux.Vars(r)
	moduleIDStr := vars["id"]

	moduleID, err := strconv.Atoi(moduleIDStr)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetPhraseModuleExercisesHandler", err, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, "invalid module ID format")
		return
	}

	if moduleID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "module ID must be positive")
		return
	}

	id := r.Context().Value(middleware.CookieName)
	userId, ok := id.(string)
	if !ok {
		userId = ""
	}

	gotModules, err := h.ucWord.GetPhraseModuleExercises(r.Context(), userId, moduleID)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetPhraseModuleExercisesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils.WriteResponse(w, http.StatusCreated, gotModules); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetPhraseModuleExercisesHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetPhraseModuleExercisesHandler")
	return
}

func (h *WordHandler) UploadTipHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "UploadTip", "parsed multipart form")

	audio_file, head_audio, err := r.FormFile("tip_audio")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audio_file.Close()

	media_file, head_media, err := r.FormFile("tip_media")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audio_file.Close()

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head_audio.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}
	phonema := r.FormValue("phonema")
	if phonema == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", errors.New("no sound"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
	}
	tip := r.FormValue("tip")
	if tip == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", errors.New("no tip"))
		utils.WriteError(w, http.StatusBadRequest, "bad data request")
	}

	tipAudioLink, err := h.minClient.UploadFile(audio_file, head_audio.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	tipPicLink, err := h.minClient.UploadFile(media_file, head_media.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}
	data := models.TipData{
		Phonema:      phonema,
		TipText:      tip,
		TipMediaLink: tipPicLink,
		TipAudioLink: tipAudioLink,
	}

	if err := h.ucWord.UploadTip(r.Context(), &data); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to upload tip")
		return
	}
	if err := utils.WriteResponse(w, http.StatusOK, "uploaded tip"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UploadTip", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}
}

func (h *WordHandler) GetTipHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())
	tip := models.TipData{}

	if err := utils.ReadRequestData(r, &tip); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTip", err, http.StatusBadRequest)
		utils.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}
	gotTip, err := h.ucWord.GetTip(r.Context(), &tip)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error get tip")
		return
	}

	audioLink, err := h.minClient.GetFileLink(gotTip.TipAudioLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetTip", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}

	gotTip.TipAudioLink = audioLink

	picLink, err := h.minClient.GetFileLink(gotTip.TipMediaLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetTip", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}

	gotTip.TipMediaLink = picLink
	if err := utils.WriteResponse(w, http.StatusOK, gotTip); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTip", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetTip")

}

func (h *WordHandler) GetCurrentModuleWordHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())
	id := r.Context().Value(middleware.CookieName)
	userId, ok := id.(string)
	if !ok {
		mod1 := models.ModuleCreate{ID: 1}
		if err := utils.WriteResponse(w, http.StatusOK, mod1); err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetCurrentModuleWordHandler", err, http.StatusInternalServerError)
			utils.WriteError(w, http.StatusInternalServerError, "error writing response")
			return
		}
		return
	}

	gotTopic, err := h.ucWord.GetNextWordModule(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error get topic progress")
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, gotTopic); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler")
}

func (h *WordHandler) GetCurrentModulePhraseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils.GetRequestIDFromCtx(r.Context())

	id := r.Context().Value(middleware.CookieName)
	userId, ok := id.(string)

	if !ok {
		mod1 := models.ModuleCreate{ID: 1}
		if err := utils.WriteResponse(w, http.StatusOK, mod1); err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetCurrentModuleWordHandler", err, http.StatusInternalServerError)
			utils.WriteError(w, http.StatusInternalServerError, "error writing response")
			return
		}
		return
	}

	gotModule, err := h.ucWord.GetNextWordModule(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error get topic progress")
		return
	}

	if err := utils.WriteResponse(w, http.StatusOK, gotModule); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler", err, http.StatusInternalServerError)
		utils.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler")

}
