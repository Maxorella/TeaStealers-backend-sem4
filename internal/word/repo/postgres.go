package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
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

func (r *WordRepo) UploadLink(ctx context.Context, tx models.Transaction, wordLink *models.WordData) error {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadLink", "new reqId")
	}

	res := r.db.QueryRow(SelectWordSql, wordLink.Word)

	wordBase := &models.WordData{}
	var Link sql.NullString
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &Link); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			r.logger.LogError(requestId, logger.RepositoryLayer, "UploadLink", err)
			return err
		}
	}

	//if Link.Valid { } TODO сейчас будем перезаписывать ссылку, мб нужно будет поменять логику

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadLink", "got word from base: "+wordBase.Word)

	if _, err := r.db.Exec(UploadLinkSql, wordLink.AudioLink, wordLink.Word); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UploadLink", err)
		return err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadLink", "return word id")
	return nil
}

func (r *WordRepo) GetWordByWord(ctx context.Context, word string) (*models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectWordSql, word)

	gotWordData := &models.WordData{}
	var AudioLink sql.NullString

	if err := res.Scan(&gotWordData.WordID, &gotWordData.Word, &gotWordData.Transcription, &AudioLink, gotWordData.Topic); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordByWord", err)
		return &models.WordData{}, err
	}

	if AudioLink.Valid {
		gotWordData.AudioLink = AudioLink.String
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetWordByWord", "got word from base: "+gotWordData.Word)
	return gotWordData, nil
}

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

func (r *WordRepo) InsertTopic(ctx context.Context, tx models.Transaction, wordTopic string) error {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	_, err := tx.Exec(InsertTag, wordTopic)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "InsertTopic", err)
		return fmt.Errorf("failed to insert topic: %w", err)
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "InsertTopic", "topic inserted successfully: "+wordTopic)
	return nil
}

func (r *WordRepo) SelectTags(ctx context.Context, tx models.Transaction) (*models.TopicsList, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "InsertTag", "new reqId")
	}

	res, err := r.db.QueryContext(ctx, SelectTags)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectTags", err)
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer res.Close()

	var tags []string
	for res.Next() {
		var tag string
		if err := res.Scan(&tag); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "SelectTags", err)
			continue // или return nil, err в зависимости от требований
		}
		tags = append(tags, tag)
	}

	if err := res.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectTags", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "SelectTags",
		fmt.Sprintf("successfully retrieved %d tags", len(tags)))
	return &models.TopicsList{Topics: tags}, nil
}

func (r *WordRepo) SelectWordsWithTopic(ctx context.Context, tx models.Transaction, tag string) (*[]models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "SelectWordsWithTopic", "new reqId") // Исправлено название метода
	}

	rows, err := r.db.QueryContext(ctx, SelectAllWordsWithTag, tag) // Используем правильный запрос
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsWithTopic", err)
		return nil, fmt.Errorf("failed to query words with tag: %w", err)
	}
	defer rows.Close()

	var words []models.WordData
	for rows.Next() {
		var word models.WordData
		var Link sql.NullString
		var Tags sql.NullString

		if err := rows.Scan(
			&word.WordID,
			&word.Word,
			&word.Transcription,
			&Link,
			&Tags,
		); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsWithTopic", err)
			continue // или return nil, err в зависимости от требований
		}
		if Link.Valid {
			word.AudioLink = Link.String
		}
		if Tags.Valid {
			word.Topic = Tags.String
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "SelectWordsWithTopic", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "SelectWordsWithTopic",
		fmt.Sprintf("successfully retrieved %d words with tag '%s'", len(words), tag))
	return &words, nil
}

/*
func (r *WordRepo) WriteStat(ctx context.Context, stat *models.WordUserStat) error {
	if stat == nil {
		return errors.New("stat is nil")
	}

	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "new reqId")
	}

	var err error
	if stat.TotalPlus == 1 {
		_, err = r.db.ExecContext(ctx, InsertPlusBigStat, stat.Id)
	} else if stat.TotalMinus == 1 {
		_, err = r.db.ExecContext(ctx, InsertMinusBigStat, stat.Id)
	}

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "WriteStat", err)
		return fmt.Errorf("failed to write stat: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "stat written successfully")
	return nil
}

func (r *WordRepo) GetStat(ctx context.Context, tx models.Transaction, word_id int) (*models.WordUserStat, error) {

	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "new reqId")
	}

	res := r.db.QueryRow(SelectBigStat, word_id)

	wordstat := &models.WordUserStat{}

	if err := res.Scan(&wordstat.Id, &wordstat.TotalPlus, &wordstat.TotalMinus); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWord", err)
			return wordstat, err
		}
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "stat written successfully")
	return wordstat, nil
}

func (r *WordRepo) UploadTip(ctx context.Context, tx models.Transaction, data *models.TipData) error {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "new reqId")
	}

	_, err := r.db.ExecContext(ctx, CreateWordTip, data.Phonema, data.TipText, data.TipPicture, data.TipAudio)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UploadTip", err)
		return err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadTip", "tip uploaded successfully")
	return nil
}

func (r *WordRepo) GetTip(ctx context.Context, tx models.Transaction, data *models.TipData) (*models.TipData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "WriteStat", "new reqId")
	}

	res := r.db.QueryRow(SelectWordTip, data.Phonema)

	gotTip := &models.TipData{}

	if err := res.Scan(&gotTip.Phonema, &gotTip.TipText, &gotTip.TipPicture, &gotTip.TipAudio); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetTip", err)
		return gotTip, err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetTip", "got tip successfully")
	return gotTip, nil
}

*/
