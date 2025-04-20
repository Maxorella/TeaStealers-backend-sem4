package stat

import (
	"context"
)

type StatUsecase interface {
	UpdateWordStat(ctx context.Context, word string, gotTranscription string) (int, error)
}
