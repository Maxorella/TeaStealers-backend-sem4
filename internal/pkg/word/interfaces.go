package word

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
)

type WordUsecase interface {
	GetWord() string
	CreateWord(ctx context.Context, wordCreateData *models.CreateWordData) (int, error)
}
