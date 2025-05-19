package module

import (
	"context"
)

type ModuleUsecase interface {
	CreateModuleWord(ctx context.Context, name string) (int, error)
	CreateModulePhrase(ctx context.Context, name string) (int, error)
}
