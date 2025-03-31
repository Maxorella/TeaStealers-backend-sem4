package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/satori/uuid"
)

type WordRepo struct {
	db     *sql.DB
	logger logger.Logger
	// cfg    *config.Config
	// metricsC metrics.MetricsHTTP
}

func NewRepository(db *sql.DB, logger logger.Logger) *WordRepo {
	return &WordRepo{db: db, logger: logger}
}

func (r *WordRepo) CreateWord(ctx context.Context, wordCreate *models.CreateWordData) (int, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWord", "new reqId")
	}

	res := r.db.QueryRow(SelectWordSql, wordCreate.Word)

	wordBase := &models.WordData{}
	var Link sql.NullString
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &Link); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWord", err)
			return -1, err
		}
	}

	if Link.Valid {
		wordBase.Link = Link.String
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWord", "got word from base: "+wordBase.Word)

	if wordCreate.Word == wordBase.Word {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWord", errors.New("word already exists"))
		return -1, errors.New("word already exists")
	}

	var lastInsertID int
	if err := r.db.QueryRowContext(ctx, CreateWordSql, wordCreate.Word, wordCreate.Transcription).Scan(&lastInsertID); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWord", err)
		return -1, err
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWord", "return word id")
	return lastInsertID, nil
}
