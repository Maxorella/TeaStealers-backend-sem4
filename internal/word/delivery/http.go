package delivery

import (
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/stat"
	"github.com/TeaStealers-backend-sem4/internal/word"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/gorilla/mux"
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
	// uc represents the usecase interface for authentication. u
	ucWord    word.WordUsecase
	ucStat    stat.StatUsecase
	cfg       *config.Config
	logger    logger.Logger
	minClient *utils2.FileStorageClient
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewWordHandler(ucWord word.WordUsecase, ucStat stat.StatUsecase, cfg *config.Config, logr logger.Logger, minCl *utils2.FileStorageClient) *WordHandler {
	return &WordHandler{ucWord: ucWord, ucStat: ucStat, cfg: cfg, logger: logr, minClient: minCl}

}

// audio.AudioUsecase GetTranscription

func (h *WordHandler) GetWord(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	vars := mux.Vars(r)
	reqWord := vars["word"]
	if reqWord == "" {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, GetWord, errors.New("empty word"), 400)
		utils2.WriteError(w, http.StatusBadRequest, "word parameter is required")
		return
	}
	wordU := &models.WordData{Word: reqWord}
	wordU.Sanitize()
	gotWord, err := h.ucWord.GetWord(r.Context(), wordU)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get word")
		return
	}

	audioLink, err := h.minClient.GetFileLink(gotWord.AudioLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetWord", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}
	gotWord.AudioLink = audioLink

	if err := utils2.WriteResponse(w, http.StatusOK, gotWord); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetWord", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetWord")

}

func (h *WordHandler) CreateWordExerciseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
		utils2.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "CreateWordExercise", "parsed multipart form")

	exercise := r.FormValue("exercise")
	if exercise == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	moduleIdStr := r.FormValue("module_id")
	if moduleIdStr == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	moduleId, err := strconv.Atoi(moduleIdStr)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("module_id not int"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
	}
	words := r.FormValue("words")
	if words == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	wordsList := utils2.ParseStringArray(words)

	transcriptions := r.FormValue("transcriptions")
	if transcriptions == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	transcriptionsList := utils2.ParseStringArray(transcriptions)

	translations := r.FormValue("translations")
	if translations == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	translationsList := utils2.ParseStringArray(translations)
	gotId := models.IdStruct{}

	switch exercise {
	case "pronounce":
		audioFile, audioHead, err := r.FormFile("audio")
		if err != nil {
			h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
			utils2.WriteError(w, http.StatusBadRequest, "bad data request")
			return
		}
		defer audioFile.Close()
		allowedExtensions := []string{".wav", ".mp3"}
		fileType := strings.ToLower(filepath.Ext(audioHead.Filename))
		if !slices.Contains(allowedExtensions, fileType) {
			utils2.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
			return
		}

		audioLink, err := h.minClient.UploadFile(audioFile, audioHead.Filename)
		if err != nil {
			h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
			utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
			return
		}

		wordData := models.CreateWordData{Exercise: exercise, ModuleId: &moduleId, Word: wordsList[0], Transcription: transcriptionsList[0], Translation: translationsList[0], AudioLink: audioLink}

		id, err := h.ucWord.CreateWordExercise(r.Context(), &wordData)
		if err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
			utils2.WriteError(w, http.StatusInternalServerError, "error create word")
			return
		}
		gotId.Id = &id
	case "pronounceFiew":
		fallthrough
	case "guessWord":
		audioFiles := r.MultipartForm.File["audio"]

		if len(audioFiles) != 2 {
			utils2.WriteError(w, http.StatusBadRequest, "not 2 audio uploaded")
			return
		}

		var audioLinks []string
		allowedExtensions := []string{".wav", ".mp3"}

		for _, fileHeader := range audioFiles {
			file, err := fileHeader.Open()
			if err != nil {
				h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
				utils2.WriteError(w, http.StatusInternalServerError, "error opening file")
				return
			}
			defer file.Close()

			fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !slices.Contains(allowedExtensions, fileExt) {
				utils2.WriteError(w, http.StatusBadRequest, "only .wav and .mp3 allowed")
				return
			}

			audioLink, err := h.minClient.UploadFile(file, fileHeader.Filename)
			if err != nil {
				h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordExercise", err)
				utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
				return
			}
			audioLinks = append(audioLinks, audioLink)
		}

		wordData := models.CreateWordDataList{Exercise: exercise, ModuleId: &moduleId, Word: wordsList,
			Transcription: transcriptionsList, Translation: translationsList, AudioLink: audioLinks}

		id, err := h.ucWord.CreateWordExerciseList(r.Context(), &wordData)
		if err != nil {
			h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
			utils2.WriteError(w, http.StatusInternalServerError, "error create word")
			return
		}
		gotId.Id = &id
	default:
		utils2.WriteError(w, http.StatusBadRequest, "no such exercise")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusCreated, gotId); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordExercise", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateWordExercise")
	return
}

func (h *WordHandler) CreatePhraseExerciseHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils2.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", "parsed multipart form")

	exercise := r.FormValue("exercise")
	if exercise == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	moduleIdStr := r.FormValue("module_id")
	if moduleIdStr == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	moduleId, err := strconv.Atoi(moduleIdStr)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("module_id not int"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
	}

	sentence := r.FormValue("sentence")
	if sentence == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	transcription := r.FormValue("transcription")
	if transcription == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	translate := r.FormValue("translate")
	if translate == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	chain := r.FormValue("chain")
	if chain == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", errors.New("bad formValue"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	chainList := utils2.ParseStringArray(chain)

	audioFile, audioHead, err := r.FormFile("audio")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audioFile.Close()
	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(audioHead.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils2.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	audioLink, err := h.minClient.UploadFile(audioFile, audioHead.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}
	gotId := models.IdStruct{}

	phraseData := models.CreatePhraseData{Exercise: exercise, Sentence: sentence, Transcription: transcription,
		ModuleId:  &moduleId,
		AudioLink: audioLink, Translate: translate, Chain: chainList}

	id, err := h.ucWord.CreatePhraseExercise(r.Context(), &phraseData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}
	gotId.Id = &id

	if err := utils2.WriteResponse(w, http.StatusCreated, gotId); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreatePhraseExerciseHandler")
	return
}

func (h *WordHandler) UpdateProgressHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())
	userId := 1 //TODO user ID !!!!!!
	progressData := models.ExerciseProgress{UserID: &userId}

	if err := utils2.ReadRequestData(r, &progressData); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	_, err := h.ucWord.CreateUpdateProgress(r.Context(), &progressData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}
	if err := utils2.WriteResponse(w, http.StatusCreated, "Progress saved"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "UpdateProgressHandler")
	return
}

/*

func (h *WordHandler) CreateWord(w http.ResponseWriter, r *http.Request) {

	requestId := utils2.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", err)
		utils2.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "CreateWordHandler", "parsed multipart form")

	audioFile, audioHead, err := r.FormFile("audio")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", err)
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audioFile.Close()
	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(audioHead.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils2.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}

	formWord := r.FormValue("word")
	if formWord == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", errors.New("no word"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	formTranscription := r.FormValue("transcription")
	if formTranscription == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", errors.New("no transcription"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	formTopic := r.FormValue("topic")
	if formTopic == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", errors.New("no topic"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}

	audioLink, err := h.minClient.UploadFile(audioFile, audioHead.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "CreateWordHandler", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	wordData := models.CreateWordData{Word: formWord, Transcription: formTranscription, Topic: formTopic, AudioLink: audioLink}

	_, err = h.ucWord.CreateWord(r.Context(), &wordData)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusCreated, "word created"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "CreateWordHandler")
	return
}
*/

func (h *WordHandler) WordsWithTopicHandler(w http.ResponseWriter, r *http.Request) {

	requestId := utils2.GetRequestIDFromCtx(r.Context())
	topic := &models.OneTopic{}

	if err := utils2.ReadRequestData(r, &topic); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", "topic to get "+topic.Topic)

	gotWords, err := h.ucStat.WordsWithTopic(r.Context(), topic.Topic)
	if err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error create word")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusCreated, gotWords); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "WordsWithTopicHandler")
	return
}

func (h *WordHandler) UploadTipHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils2.WriteError(w, http.StatusBadRequest, "max size file 5 mb")
		return
	}
	h.logger.LogInfo(requestId, logger.DeliveryLayer, "UploadTip", "parsed multipart form")

	audio_file, head_audio, err := r.FormFile("tip_audio")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audio_file.Close()

	media_file, head_media, err := r.FormFile("tip_media")
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
		return
	}
	defer audio_file.Close()

	allowedExtensions := []string{".wav", ".mp3"}
	fileType := strings.ToLower(filepath.Ext(head_audio.Filename))
	if !slices.Contains(allowedExtensions, fileType) {
		utils2.WriteError(w, http.StatusBadRequest, "wav and mp3 only")
		return
	}
	phonema := r.FormValue("phonema")
	if phonema == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", errors.New("no sound"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
	}
	tip := r.FormValue("tip")
	if tip == "" {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", errors.New("no tip"))
		utils2.WriteError(w, http.StatusBadRequest, "bad data request")
	}

	tipAudioLink, err := h.minClient.UploadFile(audio_file, head_audio.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	tipPicLink, err := h.minClient.UploadFile(media_file, head_media.Filename)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}
	data := models.TipData{
		Phonema:      phonema,
		TipText:      tip,
		TipMediaLink: tipPicLink,
		TipAudioLink: tipAudioLink,
	}
	data.Sanitize()

	if err := h.ucWord.UploadTip(r.Context(), &data); err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "UploadTip", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to upload tip")
		return
	}
	if err := utils2.WriteResponse(w, http.StatusOK, "uploaded tip"); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "UploadTip", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}
}

func (h *WordHandler) GetTipHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())
	tip := models.TipData{}

	if err := utils2.ReadRequestData(r, &tip); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTip", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}
	tip.Sanitize()
	gotTip, err := h.ucWord.GetTip(r.Context(), &tip)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get tip")
		return
	}

	audioLink, err := h.minClient.GetFileLink(gotTip.TipAudioLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetTip", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}

	gotTip.TipAudioLink = audioLink

	picLink, err := h.minClient.GetFileLink(gotTip.TipMediaLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetTip", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}

	gotTip.TipMediaLink = picLink
	if err := utils2.WriteResponse(w, http.StatusOK, gotTip); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTip", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetTip")

}

func (h *WordHandler) GetTopicProgressHandler(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())
	topic := &models.OneTopic{}

	if err := utils2.ReadRequestData(r, &topic); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	gotTopic, err := h.ucStat.GetTopicProgress(r.Context(), topic.Topic)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get topic progress")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusOK, gotTopic); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetTopicProgressHandler")

}

func (h *WordHandler) GetRandomWord(w http.ResponseWriter, r *http.Request) {
	requestId := utils2.GetRequestIDFromCtx(r.Context())

	topic := models.OneTopic{}
	if err := utils2.ReadRequestData(r, &topic); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	gotWord, err := h.ucWord.GetRandomWord(r.Context(), topic.Topic)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get word")
		return
	}

	audioLink, err := h.minClient.GetFileLink(gotWord.AudioLink)
	if err != nil {
		h.logger.LogError(requestId, logger.DeliveryLayer, "GetRandomWord", err)
		utils2.WriteError(w, http.StatusInternalServerError, "failed to get link")
		return
	}
	gotWord.AudioLink = audioLink

	if err := utils2.WriteResponse(w, http.StatusOK, gotWord); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "GetRandomWord", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "GetRandomWord")

}

/*
func (h *WordHandler) SelectTags(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
	}

	gotTags, err := h.uc.SelectTags(r.Context())
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get tags")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusOK, gotTags); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "SelectTags", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "SelectTags")

}
*/
/*
func (h *WordHandler) SelectWordsWithTopic(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
	}

	tag := models.OneTopic{}

	if err := utils2.ReadRequestData(r, &tag); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "CreateWordHandler", err, http.StatusBadRequest)
		utils2.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}
	fmt.Printf("got tag %s", tag.Topic)
	gotWords, err := h.uc.SelectWordsWithTopic(r.Context(), tag.Topic)
	if err != nil {
		utils2.WriteError(w, http.StatusInternalServerError, "error get tags")
		return
	}

	if err := utils2.WriteResponse(w, http.StatusOK, *gotWords); err != nil {
		h.logger.LogErrorResponse(requestId, logger.DeliveryLayer, "SelectWordsWithTopic", err, http.StatusInternalServerError)
		utils2.WriteError(w, http.StatusInternalServerError, "error writing response")
		return
	}

	h.logger.LogSuccessResponse(requestId, logger.DeliveryLayer, "SelectWordsWithTopic")

}
*/

/*

 */
