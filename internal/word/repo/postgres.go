package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
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

func (r *WordRepo) BeginTx(ctx context.Context) (models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *WordRepo) CreateWord(ctx context.Context, tx models.Transaction, wordCreate *models.CreateWordData) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	var lastInsertID int
	if err := tx.QueryRowContext(ctx, CreateWordSql, wordCreate.Word, wordCreate.Transcription, wordCreate.AudioLink, wordCreate.Topic).Scan(&lastInsertID); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWord", err)
		return 0, err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWord", "")
	return lastInsertID, nil
}

func (r *WordRepo) GetWordByWord(ctx context.Context, word string) (*models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectWordSql, word)

	gotWordData := &models.WordData{}
	gotWordData.WordID = new(int)
	var AudioLink sql.NullString
	if err := res.Scan(gotWordData.WordID, &gotWordData.Word, &gotWordData.Transcription, &AudioLink, &gotWordData.Topic); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordByWord", err)
		return &models.WordData{}, err
	}

	if AudioLink.Valid {
		gotWordData.AudioLink = AudioLink.String
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetWordByWord", "got word from base: "+gotWordData.Word)
	return gotWordData, nil
}

/*
func (r *WordRepo) GetRandomWord(ctx context.Context, tx models.Transaction) (*models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "new reqId")
	}

	res := r.db.QueryRow(SelectRandomWordSql)

	wordBase := &models.WordData{}
	var Link, Tags sql.NullString
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &Tags, &Link); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetRandomWord", err)
		return &models.WordData{}, err
	}

	if Link.Valid {
		wordBase.AudioLink = Link.String
	}
	if Tags.Valid {
		wordBase.Topic = Tags.String
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "got word from base: "+wordBase.Word)
	return wordBase, nil
}

func (r *WordRepo) GetRandomWordWithTag(ctx context.Context, tx models.Transaction, wordTag string) (*models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "new reqId")
	}
	res := r.db.QueryRow(SelectRandomWordWithTagSql, wordTag)

	wordBase := &models.WordData{}
	var Link, Tags sql.NullString
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &Tags, &Link); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetRandomWord", err)
		return &models.WordData{}, err
	}

	if Link.Valid {
		wordBase.AudioLink = Link.String
	}
	if Tags.Valid {
		wordBase.Topic = Tags.String
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "got word from base: "+wordBase.Word)
	return wordBase, nil
}

*/

func (r *WordRepo) SelectWordsByTopicWithProgress(ctx context.Context, tx models.Transaction, topic string) (*[]models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	rows, err := r.db.QueryContext(ctx, SelectWordWithProgressByTopic, topic)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsByTopicWithProgress", err)
		return nil, fmt.Errorf("failed to query words with tag: %w", err)
	}
	defer rows.Close()

	var words []models.WordData

	for rows.Next() {
		var word models.WordData

		if err := rows.Scan(
			&word.WordID,
			&word.Word,
			&word.Transcription,
			&word.AudioLink,
			&word.Topic,
			&word.Progress,
		); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsByTopicWithProgress", err)
			return nil, err
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsByTopicWithProgress", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "SelectWordsByTopicWithProgress",
		fmt.Sprintf("successfully retrieved"))
	return &words, nil
}

func (r *WordRepo) UploadTip(ctx context.Context, tx models.Transaction, data *models.TipData) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	_, err := r.db.ExecContext(ctx, InsertWordTip, data.Phonema, data.TipText, data.TipAudioLink, data.TipMediaLink)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UploadTip", err)
		return err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadTip", "tip uploaded successfully")
	return nil
}

func (r *WordRepo) GetTip(ctx context.Context, tx models.Transaction, data *models.TipData) (*models.TipData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectWordTip, data.Phonema)

	gotTip := &models.TipData{}

	if err := res.Scan(&gotTip.Phonema, &gotTip.TipText, &gotTip.TipAudioLink, &gotTip.TipMediaLink); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetTip", err)
		return gotTip, err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetTip", "got tip successfully")
	return gotTip, nil
}
