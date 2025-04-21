package stat

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
)

type StatUsecase interface {
	UpdateWordStat(ctx context.Context, word string, gotTranscription string) (int, error)
	GetAllTopics(ctx context.Context) (*models.TopicsList, error)
	WordsWithTopic(ctx context.Context, topic string) (*[]models.WordData, error)
	GetTopicProgress(ctx context.Context, topic string) (*models.OneTopic, error)
}
