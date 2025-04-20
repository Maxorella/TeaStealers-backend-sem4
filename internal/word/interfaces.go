package word

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
)

type WordUsecase interface {
	GetWord(ctx context.Context, wordData *models.WordData) (*models.WordData, error)
	CreateWord(ctx context.Context, wordCreateData *models.CreateWordData) (int, error)
	// UploadLink(ctx context.Context, wordLink *models.WordData) error
	// GetRandomWord(ctx context.Context, wordTag string) (*models.WordData, error)
	// SelectTags(ctx context.Context) (*models.TopicsList, error)
	// SelectWordsWithTopic(ctx context.Context, tag string) (*[]models.WordData, error)
	// WriteStat(ctx context.Context, stat *models.WordUserStat) error
	// GetStat(ctx context.Context, word_id int) (*models.WordUserStat, error)
	UploadTip(context.Context, *models.TipData) error
	GetTip(ctx context.Context, data *models.TipData) (*models.TipData, error)
}
