package usecase

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
	statRep "github.com/TeaStealers-backend-sem4/internal/stat/repo"
	wordRep "github.com/TeaStealers-backend-sem4/internal/word/repo"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
)

type StatUsecase struct {
	statRepo *statRep.StatRepo
	wordRepo *wordRep.WordRepo
	logr     logger.Logger
}

func NewStatUsecase(statrep *statRep.StatRepo, wordrep *wordRep.WordRepo, logger logger.Logger) *StatUsecase {
	return &StatUsecase{statRepo: statrep, wordRepo: wordrep, logr: logger}
}

func (uc *StatUsecase) UpdateWordStat(ctx context.Context, word string, gotTranscription string) (int, error) {
	//	requestId := utils2.GetRequestIDFromCtx(ctx)
	gotWord, err := uc.wordRepo.GetWordByWord(ctx, word)
	if err != nil {
		return 0, err //TODO errors
	} // TODO error handling usecase
	tx, err := uc.statRepo.BeginTx(ctx)

	if err != nil {
		return 0, err
	}
	var result int

	if gotTranscription == gotWord.Transcription {
		result = 1
	} else {
		result = -1
	}

	err = uc.statRepo.UpdateWordStat(ctx, tx, gotWord, result)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return result, nil
}

func (uc *StatUsecase) GetAllTopics(ctx context.Context) (*models.TopicsList, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.statRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	allTopics, err := uc.statRepo.SelectAllTopics(ctx, tx)
	if err != nil {
		uc.logr.LogError(requestId, logger.UsecaseLayer, "GetAllTopics", err)
		return nil, err
	}
	for i := range allTopics.Topics {
		newtopic, err := uc.statRepo.GetTopicProgress(ctx, tx, allTopics.Topics[i].Topic)
		if err != nil {
			uc.logr.LogError(requestId, logger.UsecaseLayer, "GetAllTopics", err)
			continue
		}
		allTopics.Topics[i] = *newtopic
	}
	uc.logr.LogInfo(requestId, logger.UsecaseLayer, "GetAllTopics", "finished GetAllTopics usecase")

	return allTopics, nil
}

func (uc *StatUsecase) WordsWithTopic(ctx context.Context, topic string) (*[]models.WordData, error) {
	// requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.statRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	gotWords, err := uc.wordRepo.SelectWordsByTopicWithProgress(ctx, tx, topic)
	if err != nil {
		return nil, err
	}
	return gotWords, nil
}

func (uc *StatUsecase) GetTopicProgress(ctx context.Context, topic string) (*models.OneTopic, error) {
	// requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.statRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	gotTopic, err := uc.statRepo.GetTopicProgress(ctx, tx, topic)
	if err != nil {
		return nil, err
	}

	return gotTopic, nil
}
