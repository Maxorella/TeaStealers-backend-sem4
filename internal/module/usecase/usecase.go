package usecase

import (
	"context"
	moduleRep "github.com/TeaStealers-backend-sem4/internal/module/repo"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
)

type ModuleUsecase struct {
	Repo *moduleRep.ModuleRepo
	logr logger.Logger
}

func NewModuleUsecase(statrep *moduleRep.ModuleRepo, logger logger.Logger) *ModuleUsecase {
	return &ModuleUsecase{Repo: statrep, logr: logger}
}

func (uc *ModuleUsecase) CreateModuleWord(ctx context.Context, moduleName string) (int, error) {
	//	requestId := utils.GetRequestIDFromCtx(ctx)
	tx, err := uc.Repo.BeginTx(ctx)
	if err != nil {
		return -1, err
	}

	gotId, err := uc.Repo.InsertModuleWord(ctx, tx, moduleName)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return gotId, nil
}

func (uc *ModuleUsecase) CreateModulePhrase(ctx context.Context, moduleName string) (int, error) {
	//	requestId := utils.GetRequestIDFromCtx(ctx)
	tx, err := uc.Repo.BeginTx(ctx)
	if err != nil {
		return -1, err
	}

	gotId, err := uc.Repo.InsertModulePhrase(ctx, tx, moduleName)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return gotId, nil
}
