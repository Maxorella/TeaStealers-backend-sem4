package audio

import (
	"context"
)

type AudioUsecase interface {
	GetTranscription(ctx context.Context, audio []byte) (string, error)
}
