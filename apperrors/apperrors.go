package apperrors

import "errors"

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
