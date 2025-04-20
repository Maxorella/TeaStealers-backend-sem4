package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"strconv"
)

type StatRepo struct {
	db     *sql.DB
	logger logger.Logger
	// cfg    *config.Config
	// metricsC metrics.MetricsHTTP
}

func NewRepository(db *sql.DB, logger logger.Logger) *StatRepo {
	return &StatRepo{db: db, logger: logger}
}

func (r *StatRepo) BeginTx(ctx context.Context) (models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *StatRepo) UpdateWordStat(ctx context.Context, tx models.Transaction, word *models.WordData, result int) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	wordIdStr := strconv.Itoa(*word.WordID)

	res, err := tx.ExecContext(ctx, CreateUpdateWordStatSql, wordIdStr, result)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to update word stat: %w", err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		errMsg := "no rows were affected by the update"
		r.logger.LogError(requestId, logger.RepositoryLayer, "UpdateWordStat", errors.New(errMsg))
		return errors.New(errMsg)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UpdateWordStat",
		fmt.Sprintf("successfully updated word stat, word_id: %s, rows affected: %d", wordIdStr, rowsAffected))
	return nil
}
