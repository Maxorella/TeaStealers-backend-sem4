package usecase

import (
	"context"
	"errors"
)

type AudioUsecase struct {
}

func NewAudioUsecase() *AudioUsecase {
	return &AudioUsecase{}
}

// SignUp handles the user registration process.
func (u *AudioUsecase) GetTranscription(ctx context.Context, audio []byte) (string, error) {
	return "", errors.New("")
}
