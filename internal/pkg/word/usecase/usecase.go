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

func (uc *WordUsecase) CreateWord(ctx context.Context, wordCreateData *models.CreateWordData) (int, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "new reqId")
	}

	word_id, err := uc.repo.CreateWord(ctx, wordCreateData)
	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return -1, errors.New("failed to create word")
	}
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "created new word uc")
	return word_id, nil
}

func (uc *WordUsecase) UploadLink(ctx context.Context, wordLink *models.WordData) error {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "new reqId")
	}

	err := uc.repo.UploadLink(ctx, wordLink)

	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return errors.New("failed to upload uuid")
	}
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", "uploaded uuid to word")
	return nil
}

func (uc *WordUsecase) GetWord(ctx context.Context, wordData *models.WordData) (*models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetWord", "new reqId")
	}

	gotWord, err := uc.repo.GetWord(ctx, wordData)

	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetWord", err)
		return gotWord, errors.New("failed to get word")
	}
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetWord", "uploaded uuid to word")
	return gotWord, nil
}
