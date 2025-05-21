package word

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
)

type WordUsecase interface {
	CreateWordExercise(ctx context.Context, wordCreateData *models.CreateWordData) (int, error)
	CreateWordExerciseList(ctx context.Context, wordCreateData *models.CreateWordDataList) (int, error)
	CreatePhraseExercise(ctx context.Context, phraseCreateData *models.CreatePhraseData) (int, error)

	CreateUpdateProgress(ctx context.Context, progress *models.ExerciseProgress) (int, error)

	GetWordModuleExercises(ctx context.Context, userID string, moduleId int) (*models.ExerciseList, error)
	GetPhraseModuleExercises(ctx context.Context, userID string, moduleId int) (*models.ExerciseList, error)

	GetWordModules(ctx context.Context) (*models.ModuleList, error)
	GetPhraseModules(ctx context.Context) (*models.ModuleList, error)

	GetNextPhraseModule(ctx context.Context, userID string) (*models.ModuleCreate, error)
	GetNextWordModule(ctx context.Context, userID string) (*models.ModuleCreate, error)

	UploadTip(ctx context.Context, data *models.TipData) error
	GetTip(ctx context.Context, data *models.TipData) (*models.TipData, error)
}
