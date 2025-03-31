package usecase

import (
	"context"
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/word/repo"
	"github.com/satori/uuid"
)

type WordUsecase struct {
	repo   *repo.WordRepo
	logger logger.Logger
}

func NewWordUsecase(repo *repo.WordRepo, logger logger.Logger) *WordUsecase {
	return &WordUsecase{
		repo:   repo,
		logger: logger,
	}
}

// SignUp handles the user registration process.
func (uc *WordUsecase) GetWord() string {
	return ""
}

func (uc *WordUsecase) CreateWord(ctx context.Context, wordCreateData *models.CreateWordData) (int, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "new reqId")
	}

	word_id, err := uc.repo.CreateWord(ctx, wordCreateData)
	if err != nil {
		// TODO logs
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return -1, errors.New("failed to create word")
	}
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "created new word uc")
	return word_id, nil
}
