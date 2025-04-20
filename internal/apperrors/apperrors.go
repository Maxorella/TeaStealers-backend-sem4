package apperrors

import (
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	servererrors "github.com/TeaStealers-backend-sem4/pkg/server_errors"
	"net/http"
)

var errToCode = map[error]int{
	// HTTP errors
	servererrors.ErrInvalidRequestData: http.StatusBadRequest, // 400

	// Usecase errors

	// Repository errors

}

// Ошибки, связанные с параметрами.
var (
	ErrInvalidParams            = errors.New("invalid parameters")
	ErrObjectIDRequired         = errors.New("objectID required")
	ErrFileMultipartKeyRequired = errors.New("bad data request file parameter required")
)

// Ошибки, связанные с аудиофайлами
var (
	ErrMaxFileSize5 = errors.New("max size file 5 mb")
	ErrWavMp3Only   = errors.New("bad request, wav and mp3 only")
	ErrReadFileForm = errors.New("unable to read file form")
)

// Ошибки, связанные с клиентом Minio
var (
	ErrFailedToGetFIle = errors.New("failed to get file")
	ErrFailedSaveFile  = errors.New("failed to save file")
)

// Ошибки, связанные с сервером.
var (
	ErrInternalServer = errors.New("internal server error")
)

// Ошибки в usecase
var (
	ErrCreateWOrd = errors.New("failed to create word")
)

func GetErrAndCodeToSend(err error) (int, error) {
	var source error
	for err != nil {
		if errors.Is(err, models.ErrValidation) {
			return http.StatusBadRequest, err
		}
		source = err
		err = errors.Unwrap(err)
	}

	code, ok := errToCode[source]
	if !ok {
		return http.StatusInternalServerError, servererrors.ErrInternalServer
	}
	return code, source
}
