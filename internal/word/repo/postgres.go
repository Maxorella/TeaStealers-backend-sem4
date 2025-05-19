package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/lib/pq"
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

func (r *WordRepo) CreateWordExercise(ctx context.Context, tx models.Transaction, wordCreate *models.CreateWordData) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	// Преобразуем массивы в PostgreSQL массивы
	wordsArray := pq.Array([]string{wordCreate.Word})
	transcriptionsArray := pq.Array([]string{wordCreate.Transcription})
	translationsArray := pq.Array([]string{wordCreate.Translation})
	audioArray := pq.Array([]string{wordCreate.AudioLink}) // Создаем массив с одним элементом

	var lastInsertID int
	err := tx.QueryRowContext(ctx, CreateWordExerciseSql,
		wordCreate.Exercise, // Тип упражнения
		wordsArray,
		transcriptionsArray,
		audioArray,
		translationsArray,
		wordCreate.ModuleId,
	).Scan(&lastInsertID)

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWordExercise", err)
		return 0, fmt.Errorf("failed to create word exercise: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWordExercise", "word exercise created")
	return lastInsertID, nil
}

func (r *WordRepo) CreatePhraseExercise(ctx context.Context, tx models.Transaction, phraseCreate *models.CreatePhraseData) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	// Валидация типа упражнения
	validExercises := map[string]bool{
		"pronounce":     true,
		"completeChain": true,
	}

	if !validExercises[phraseCreate.Exercise] {
		return 0, fmt.Errorf("invalid exercise type: %s", phraseCreate.Exercise)
	}

	// Для упражнения "completeChain" требуется цепочка слов
	if phraseCreate.Exercise == "completeChain" && len(phraseCreate.Chain) == 0 {
		return 0, errors.New("chain exercise requires at least one word in chain")
	}

	var lastInsertID int
	err := tx.QueryRowContext(ctx, CreatePhraseExerciseSql,
		phraseCreate.Exercise,
		phraseCreate.Sentence,
		phraseCreate.Translate,
		phraseCreate.Transcription,
		phraseCreate.AudioLink,
		pq.Array(phraseCreate.Chain),
		phraseCreate.ModuleId,
	).Scan(&lastInsertID)

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreatePhraseExercise", err)
		return 0, fmt.Errorf("failed to create phrase exercise: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreatePhraseExercise", "phrase exercise created")
	return lastInsertID, nil
}

func (r *WordRepo) CreateOrUpdateExerciseProgress(
	ctx context.Context,
	tx models.Transaction,
	progress *models.ExerciseProgress,
) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	// Валидация типа упражнения
	validTypes := map[string]bool{
		"word":   true,
		"phrase": true,
	}
	if !validTypes[progress.ExerciseType] {
		return 0, fmt.Errorf("invalid exercise type: %s", progress.ExerciseType)
	}

	// Валидация статуса
	validStatuses := map[string]bool{
		"none":        true,
		"in_progress": true,
		"completed":   true,
		"failed":      true,
	}
	if !validStatuses[progress.Status] {
		return 0, fmt.Errorf("invalid status: %s", progress.Status)
	}

	var lastInsertID int
	err := tx.QueryRowContext(
		ctx,
		UpsertExerciseProgressSql,
		progress.UserID,
		progress.ExerciseID,
		progress.ExerciseType,
		progress.Status,
	).Scan(&lastInsertID)

	if err != nil {
		r.logger.LogError(
			requestId,
			logger.RepositoryLayer,
			"CreateOrUpdateExerciseProgress",
			err,
		)
		return 0, fmt.Errorf("failed to upsert exercise progress: %w", err)
	}

	r.logger.LogInfo(
		requestId,
		logger.RepositoryLayer,
		"CreateOrUpdateExerciseProgress",
		"exercise progress upserted successfully",
	)
	return lastInsertID, nil
}

func (r *WordRepo) CreateWordExerciseList(ctx context.Context, tx models.Transaction, wordCreate *models.CreateWordDataList) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	// Преобразуем массивы в PostgreSQL массивы
	wordsArray := pq.Array(wordCreate.Word)
	transcriptionsArray := pq.Array(wordCreate.Transcription)
	translationsArray := pq.Array(wordCreate.Translation)
	audioArray := pq.Array(wordCreate.AudioLink) // Создаем массив с одним элементом

	var lastInsertID int
	err := tx.QueryRowContext(ctx, CreateWordExerciseSql,
		wordCreate.Exercise, // Тип упражнения
		wordsArray,
		transcriptionsArray,
		audioArray,
		translationsArray,
		wordCreate.ModuleId,
	).Scan(&lastInsertID)

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "CreateWordExercise", err)
		return 0, fmt.Errorf("failed to create word exercise: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "CreateWordExercise", "word exercise created")
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

func (r *WordRepo) GetRandomWord(ctx context.Context, tx models.Transaction) (*models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectRandomWordSql)

	wordBase := &models.WordData{}
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &wordBase.Topic, &wordBase.AudioLink); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetRandomWord", err)
		return &models.WordData{}, err
	}
	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "got word from base: "+wordBase.Word)
	return wordBase, nil
}

func (r *WordRepo) GetRandomWordWithTag(ctx context.Context, tx models.Transaction, wordTopic string) (*models.WordData, error) {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = uuid.NewV4().String()
		ctx = context.WithValue(ctx, "requestId", requestId)
		r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "new reqId")
	}
	res := r.db.QueryRow(SelectRandomWordWithTopicSql, wordTopic)

	wordBase := &models.WordData{}
	if err := res.Scan(&wordBase.WordID, &wordBase.Word, &wordBase.Transcription, &wordBase.Topic, &wordBase.AudioLink); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetRandomWord", err)
		return &models.WordData{}, err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetRandomWord", "got word from base: "+wordBase.Word)
	return wordBase, nil
}

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
			// &word.WordID,
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
