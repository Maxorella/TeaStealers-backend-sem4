package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/word/repo"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
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
	requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.repo.BeginTx(ctx)

	wordId, err := uc.repo.CreateWord(ctx, tx, wordCreateData)
	if err != nil {
		tx.Rollback()
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return 0, errors.New("failed to create word")
	}

	err = uc.repo.InsertTopic(ctx, tx, wordCreateData.Topic)

	if err != nil {
		tx.Rollback()
		return 0, errors.New("failed to create topic")
	}
	tx.Commit()
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", fmt.Sprintf("created new word, id: %d", word_id))
	return wordId, nil
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
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetWord", "got word")
	return gotWord, nil
}

func (uc *WordUsecase) GetRandomWord(ctx context.Context, wordTag string) (*models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetRandomWord", "new reqId")
	}
	var gotWord *models.WordData
	var err error
	if wordTag == "" {
		gotWord, err = uc.repo.GetRandomWord(ctx)
	} else {
		gotWord, err = uc.repo.GetRandomWordWithTag(ctx, wordTag)

	}
	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetRandomWord", err)
		return gotWord, errors.New("failed to get word")
	}

	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetRandomWord", "got word")
	return gotWord, nil
}

func (uc *WordUsecase) SelectTags(ctx context.Context) (*models.TagsList, error) {
	gottags, err := uc.repo.SelectTags(ctx)
	return gottags, err
}

func (uc *WordUsecase) SelectWordsWithTopic(ctx context.Context, tag string) (*[]models.WordData, error) {
	gotwords, err := uc.repo.SelectWordsWithTopic(ctx, tag)
	return gotwords, err
}

func (uc *WordUsecase) WriteStat(ctx context.Context, stat *models.WordStat) error {
	err := uc.repo.WriteStat(ctx, stat)
	return err
}

func (uc *WordUsecase) GetStat(ctx context.Context, word_id int) (*models.WordStat, error) {
	stat, err := uc.repo.GetStat(ctx, word_id)
	return stat, err
}

func (uc *WordUsecase) UploadTip(ctx context.Context, data *models.TipData) error {
	err := uc.repo.UploadTip(ctx, data)
	return err
}

func (uc *WordUsecase) GetTip(ctx context.Context, data *models.TipData) (*models.TipData, error) {
	gotTip, err := uc.repo.GetTip(ctx, data)
	return gotTip, err
}
