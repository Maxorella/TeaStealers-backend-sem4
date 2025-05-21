package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/lib/pq"
)

type WordRepo struct {
	db     *sql.DB
	logger logger.Logger
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
	requestId := utils.GetRequestIDFromCtx(ctx)

	wordsArray := pq.Array([]string{wordCreate.Word})
	transcriptionsArray := pq.Array([]string{wordCreate.Transcription})
	translationsArray := pq.Array([]string{wordCreate.Translation})
	audioArray := pq.Array([]string{wordCreate.AudioLink})

	var lastInsertID int
	err := tx.QueryRowContext(ctx, CreateWordExerciseSql,
		wordCreate.Exercise,
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
	requestId := utils.GetRequestIDFromCtx(ctx)

	validExercises := map[string]bool{
		"pronounce":     true,
		"completeChain": true,
	}

	if !validExercises[phraseCreate.Exercise] {
		return 0, fmt.Errorf("invalid exercise type: %s", phraseCreate.Exercise)
	}

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
	requestId := utils.GetRequestIDFromCtx(ctx)

	validTypes := map[string]bool{
		"word":   true,
		"phrase": true,
	}
	if !validTypes[progress.ExerciseType] {
		return 0, fmt.Errorf("invalid exercise type: %s", progress.ExerciseType)
	}

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
	requestId := utils.GetRequestIDFromCtx(ctx)

	wordsArray := pq.Array(wordCreate.Word)
	transcriptionsArray := pq.Array(wordCreate.Transcription)
	translationsArray := pq.Array(wordCreate.Translation)
	audioArray := pq.Array(wordCreate.AudioLink)

	var lastInsertID int
	err := tx.QueryRowContext(ctx, CreateWordExerciseSql,
		wordCreate.Exercise,
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

func (r *WordRepo) GetPhraseModules(ctx context.Context) (*models.ModuleList, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	rows, err := r.db.QueryContext(ctx, SelectPhraseModulesSql)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModules", err)
		return nil, fmt.Errorf("failed to get phrase modules: %w", err)
	}
	defer rows.Close()

	var modules []models.ModuleCreate
	for rows.Next() {
		var module models.ModuleCreate
		if err := rows.Scan(&module.ID, &module.Title); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModules", err)
			return nil, fmt.Errorf("failed to scan phrase module: %w", err)
		}
		modules = append(modules, module)
	}

	if err = rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModules", err)
		return nil, fmt.Errorf("error after iterating phrase modules: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetPhraseModules",
		fmt.Sprintf("retrieved %d phrase modules", len(modules)))

	return &models.ModuleList{Modules: modules}, nil
}

func (r *WordRepo) GetWordModules(ctx context.Context) (*models.ModuleList, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	rows, err := r.db.QueryContext(ctx, SelectWordModulesSql)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModules", err)
		return nil, fmt.Errorf("failed to get word modules: %w", err)
	}
	defer rows.Close()

	var modules []models.ModuleCreate
	for rows.Next() {
		var module models.ModuleCreate
		if err := rows.Scan(&module.ID, &module.Title); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModules", err)
			return nil, fmt.Errorf("failed to scan word module: %w", err)
		}
		modules = append(modules, module)
	}

	if err = rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModules", err)
		return nil, fmt.Errorf("error after iterating word modules: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetWordModules",
		fmt.Sprintf("retrieved %d word modules", len(modules)))

	return &models.ModuleList{Modules: modules}, nil
}

func (r *WordRepo) UploadTip(ctx context.Context, tx models.Transaction, data *models.TipData) error {
	requestId := utils.GetRequestIDFromCtx(ctx)

	_, err := r.db.ExecContext(ctx, InsertWordTip, data.Phonema, data.TipText, data.TipAudioLink, data.TipMediaLink)
	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "UploadTip", err)
		return err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "UploadTip", "tip uploaded successfully")
	return nil
}

func (r *WordRepo) GetTip(ctx context.Context, tx models.Transaction, data *models.TipData) (*models.TipData, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	res := r.db.QueryRow(SelectWordTip, data.Phonema)

	gotTip := &models.TipData{}

	if err := res.Scan(&gotTip.Phonema, &gotTip.TipText, &gotTip.TipAudioLink, &gotTip.TipMediaLink); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetTip", err)
		return gotTip, err
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetTip", "got tip successfully")
	return gotTip, nil
}

func (r *WordRepo) GetWordModuleExercises(ctx context.Context, userID string, moduleID int) (*models.ExerciseList, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	var rows *sql.Rows
	var err error

	if len(userID) > 0 {
		rows, err = r.db.QueryContext(ctx, GetWordModuleExercisesWithProgressSql, userID, moduleID)
	} else {
		rows, err = r.db.QueryContext(ctx, GetWordModuleExercisesSql, moduleID)
	}

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModuleExercises", err)
		return nil, fmt.Errorf("failed to query word exercises: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var exercise models.Exercise
		var words, transcriptions, audio, translations pq.StringArray

		if err := rows.Scan(
			&exercise.ID,
			&exercise.ExerciseType,
			&words,
			&transcriptions,
			&audio,
			&translations,
			&exercise.ModuleId,
			&exercise.Status,
		); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModuleExercises", err)
			return nil, fmt.Errorf("failed to scan word exercise: %w", err)
		}

		exercise.Words = words
		exercise.Transcriptions = transcriptions
		exercise.Audio = audio
		exercise.Translations = translations

		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetWordModuleExercises", err)
		return nil, fmt.Errorf("error after iterating word exercises: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetWordModuleExercises",
		fmt.Sprintf("retrieved %d exercises for word module %d", len(exercises), moduleID))

	return &models.ExerciseList{Exercises: exercises}, nil
}

func (r *WordRepo) GetPhraseModuleExercises(ctx context.Context, userID string, moduleID int) (*models.ExerciseList, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	var rows *sql.Rows
	var err error

	if len(userID) > 0 {
		rows, err = r.db.QueryContext(ctx, GetPhraseModuleExercisesWithProgressSql, userID, moduleID)
	} else {
		rows, err = r.db.QueryContext(ctx, GetPhraseModuleExercisesSql, moduleID)
	}

	if err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModuleExercises", err)
		return nil, fmt.Errorf("failed to query phrase exercises: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var exercise models.Exercise
		var sentence, translate, transcription, audio string
		var chain pq.StringArray
		var status string

		if err := rows.Scan(
			&exercise.ID,
			&exercise.ExerciseType,
			&sentence,
			&translate,
			&transcription,
			&audio,
			&chain,
			&exercise.ModuleId,
			&status,
		); err != nil {
			r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModuleExercises", err)
			return nil, fmt.Errorf("failed to scan phrase exercise: %w", err)
		}

		exercise.Words = []string{sentence}
		exercise.Translations = []string{translate}
		exercise.Transcriptions = []string{transcription}
		exercise.Audio = []string{audio}
		exercise.Chain = chain
		exercise.Status = status

		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetPhraseModuleExercises", err)
		return nil, fmt.Errorf("error after iterating phrase exercises: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetPhraseModuleExercises",
		fmt.Sprintf("retrieved %d exercises for phrase module %d", len(exercises), moduleID))

	return &models.ExerciseList{Exercises: exercises}, nil
}

func (r *WordRepo) GetIncompletePhraseModule(ctx context.Context, userID string) (*models.ModuleCreate, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	var module models.ModuleCreate

	err := r.db.QueryRowContext(ctx, GetIncompletePhraseModuleSql, userID).Scan(&module.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetIncompletePhraseModule", "no incomplete modules found")
			return nil, nil
		}
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetIncompletePhraseModule", err)
		return nil, fmt.Errorf("failed to get incomplete phrase module: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetIncompletePhraseModule",
		fmt.Sprintf("found incomplete phrase module ID: %d", module.ID))

	return &module, nil
}

func (r *WordRepo) GetIncompleteWordModule(ctx context.Context, userID string) (*models.ModuleCreate, error) {
	requestId := utils.GetRequestIDFromCtx(ctx)

	var module models.ModuleCreate

	err := r.db.QueryRowContext(ctx, GetIncompleteWordModuleSql, userID).Scan(&module.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetIncompleteWordModule", "no incomplete modules found")
			return nil, nil
		}
		r.logger.LogError(requestId, logger.RepositoryLayer, "GetIncompleteWordModule", err)
		return nil, fmt.Errorf("failed to get incomplete word module: %w", err)
	}

	r.logger.LogInfo(requestId, logger.RepositoryLayer, "GetIncompleteWordModule",
		fmt.Sprintf("found incomplete word module ID: %d", module.ID))

	return &module, nil
}
