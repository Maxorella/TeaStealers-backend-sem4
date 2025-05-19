package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
)

type ModuleRepo struct {
	db     *sql.DB
	logger logger.Logger
	// cfg    *config.Config
	// metricsC metrics.MetricsHTTP
}

func NewRepository(db *sql.DB, logger logger.Logger) *ModuleRepo {
	return &ModuleRepo{db: db, logger: logger}
}

func (r *ModuleRepo) BeginTx(ctx context.Context) (models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *ModuleRepo) InsertModuleWord(ctx context.Context, tx models.Transaction, moduleName string) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	var id int
	err := r.db.QueryRow(CreateModuleWord, moduleName).Scan(&id)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "InsertModuleWord", err)
		return 0, fmt.Errorf("failed to insert word module: %w", err)
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "InsertModuleWord", "word module created")
	return id, nil
}

func (r *ModuleRepo) InsertModulePhrase(ctx context.Context, tx models.Transaction, moduleName string) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	var id int
	err := r.db.QueryRow(CreateModulePhrase, moduleName).Scan(&id)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "InsertModulePhrase", err)
		return 0, fmt.Errorf("failed to insert phrase module: %w", err)
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "InsertModulePhrase", "phrase module created")
	return id, nil
}
