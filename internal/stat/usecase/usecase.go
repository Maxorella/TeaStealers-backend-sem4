package usecase

import (
	"context"
	statRep "github.com/TeaStealers-backend-sem4/internal/stat/repo"
	wordRep "github.com/TeaStealers-backend-sem4/internal/word/repo"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
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
