package word

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
)

type WordUsecase interface {
	GetWord(ctx context.Context, wordData *models.WordData) (*models.WordData, error)
	CreateWord(ctx context.Context, wordCreateData *models.CreateWordData) (int, error)
	UploadLink(ctx context.Context, wordLink *models.WordData) error
	GetRandomWord(ctx context.Context, wordTag string) (*models.WordData, error)
}
