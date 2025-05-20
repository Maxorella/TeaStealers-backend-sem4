package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/models"
	statRep "github.com/TeaStealers-backend-sem4/internal/stat/repo"
	"github.com/TeaStealers-backend-sem4/internal/word/repo"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
)

type WordUsecase struct {
	wordRepo *repo.WordRepo
	statRepo *statRep.StatRepo
	logger   logger.Logger
}

func NewWordUsecase(repoWord *repo.WordRepo, repoStat *statRep.StatRepo, logger logger.Logger) *WordUsecase {
	return &WordUsecase{
		wordRepo: repoWord,
		statRepo: repoStat,
		logger:   logger,
	}
}

func (uc *WordUsecase) CreateWordExercise(ctx context.Context, wordCreateData *models.CreateWordData) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return 0, errors.New("error begin tx")
	}
	/*
		_, err = uc.wordRepo.GetWordByWord(ctx, wordCreateData.Word)
		if err == nil {
			uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", errors.New("word already exists"))
			return 0, errors.New("word already exists")
		}
		uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", err.Error())
	*/

	wordId, err := uc.wordRepo.CreateWordExercise(ctx, tx, wordCreateData)
	if err != nil {
		tx.Rollback()
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return 0, errors.New("failed to create word")
	}
	/*
		err = uc.statRepo.InsertTopic(ctx, tx, wordCreateData.Topic)

		if err != nil {
			tx.Rollback()
			return 0, errors.New("failed to create topic")
		}

	*/
	tx.Commit()
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", fmt.Sprintf("created new word, id: %d", wordId))
	return wordId, nil
}

func (uc *WordUsecase) CreateWordExerciseList(ctx context.Context, wordCreateData *models.CreateWordDataList) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return 0, errors.New("error begin tx")
	}

	wordId, err := uc.wordRepo.CreateWordExerciseList(ctx, tx, wordCreateData)
	if err != nil {
		tx.Rollback()
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return 0, errors.New("failed to create word")
	}
	tx.Commit()
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", fmt.Sprintf("created new word, id: %d", wordId))
	return wordId, nil
}

func (uc *WordUsecase) CreatePhraseExercise(ctx context.Context, phraseCreateData *models.CreatePhraseData) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return 0, errors.New("error begin tx")
	}

	wordId, err := uc.wordRepo.CreatePhraseExercise(ctx, tx, phraseCreateData)
	if err != nil {
		tx.Rollback()
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateWord", err)
		return 0, errors.New("failed to create word")
	}
	tx.Commit()
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateWord", fmt.Sprintf("created new word, id: %d", wordId))
	return wordId, nil
}

func (uc *WordUsecase) CreateUpdateProgress(ctx context.Context, progress *models.ExerciseProgress) (int, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)

	// Начинаем транзакцию
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateUpdateProgress",
			fmt.Errorf("failed to begin transaction: %w", err))
		return 0, errors.New("failed to start transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Валидация прогресса перед сохранением
	/*
		if  *progress.ExerciseID == 0 {
			err = errors.New("exercise_id")
			uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateUpdateProgress", err)
			return 0, err
		}

	*/

	// Сохраняем/обновляем прогресс
	progressID, err := uc.wordRepo.CreateOrUpdateExerciseProgress(ctx, tx, progress)
	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateUpdateProgress",
			fmt.Errorf("failed to save progress: %w", err))
		return 0, errors.New("failed to save progress")
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "CreateUpdateProgress",
			fmt.Errorf("failed to commit transaction: %w", err))
		return 0, errors.New("failed to save progress")
	}

	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "CreateUpdateProgress",
		fmt.Sprintf("successfully saved progress, id: %d", progressID))

	return progressID, nil
}

func (uc *WordUsecase) GetWord(ctx context.Context, wordData *models.WordData) (*models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	gotWord, err := uc.wordRepo.GetWordByWord(ctx, wordData.Word)

	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetWord", err)
		return gotWord, errors.New("failed to get word")
	}
	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetWord", "got word")
	return gotWord, nil
}

func (uc *WordUsecase) GetRandomWord(ctx context.Context, wordTopic string) (*models.WordData, error) {
	requestId := utils2.GetRequestIDFromCtx(ctx)
	var gotWord *models.WordData
	var err error
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	if wordTopic == "" {
		gotWord, err = uc.wordRepo.GetRandomWord(ctx, tx)
	} else {
		gotWord, err = uc.wordRepo.GetRandomWordWithTag(ctx, tx, wordTopic)

	}
	if err != nil {
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetRandomWord", err)
		return gotWord, errors.New("failed to get word")
	}

	uc.logger.LogInfo(requestId, logger.UsecaseLayer, "GetRandomWord", "got word")
	return gotWord, nil
}

/*
	func (uc *WordUsecase) SelectWordsWithTopic(ctx context.Context, tag string) (*[]models.WordData, error) {
		gotwords, err := uc.wordRepo.SelectWordsByTopicWithProgress(ctx, tag)
		return gotwords, err
	}
*/
func (uc *WordUsecase) GetPhraseModules(ctx context.Context) (*models.ModuleList, error) {
	modules, err := uc.wordRepo.GetPhraseModules(ctx)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetPhraseModules", err)
		return nil, fmt.Errorf("failed to get phrase modules: %w", err)
	}
	return modules, nil
}

// GetWordModules возвращает список модулей для слов
func (uc *WordUsecase) GetWordModules(ctx context.Context) (*models.ModuleList, error) {
	modules, err := uc.wordRepo.GetWordModules(ctx)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetWordModules", err)
		return nil, fmt.Errorf("failed to get word modules: %w", err)
	}
	return modules, nil
}

func (uc *WordUsecase) GetWordModuleExercises(ctx context.Context, userID string, moduleId int) (*models.ExerciseList, error) {
	modules, err := uc.wordRepo.GetWordModuleExercises(ctx, userID, moduleId)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetWordModuleExercises", err)
		return nil, fmt.Errorf("failed to  modules: %w", err)
	}
	return modules, nil
}

func (uc *WordUsecase) GetPhraseModuleExercises(ctx context.Context, userID string, moduleId int) (*models.ExerciseList, error) {
	modules, err := uc.wordRepo.GetPhraseModuleExercises(ctx, userID, moduleId)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetWordModuleExercises", err)
		return nil, fmt.Errorf("failed to  modules: %w", err)
	}
	return modules, nil
}

func (uc *WordUsecase) UploadTip(ctx context.Context, data *models.TipData) error {
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	err = uc.wordRepo.UploadTip(ctx, tx, data)
	return err
}

func (uc *WordUsecase) GetTip(ctx context.Context, data *models.TipData) (*models.TipData, error) {
	tx, err := uc.wordRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	gotTip, err := uc.wordRepo.GetTip(ctx, tx, data)
	return gotTip, err
}

func (uc *WordUsecase) GetNextPhraseModule(ctx context.Context, userID string) (*models.ModuleCreate, error) {
	module, err := uc.wordRepo.GetIncompletePhraseModule(ctx, userID)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetNextPhraseModule", err)
		return nil, fmt.Errorf("failed to get next phrase module: %w", err)
	}
	return module, nil
}

func (uc *WordUsecase) GetNextWordModule(ctx context.Context, userID string) (*models.ModuleCreate, error) {
	module, err := uc.wordRepo.GetIncompleteWordModule(ctx, userID)
	if err != nil {
		requestId := utils2.GetRequestIDFromCtx(ctx)
		uc.logger.LogError(requestId, logger.UsecaseLayer, "GetNextWordModule", err)
		return nil, fmt.Errorf("failed to get next word module: %w", err)
	}
	return module, nil
}
