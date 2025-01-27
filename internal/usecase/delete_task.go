package usecase

import (
	"context"
	"test-ns/internal/entities"
	"time"
)

type (
	DeleteTaskUseCase interface {
		Execute(ctx context.Context, input DeleteTaskInput) error
	}

	DeleteTaskInput struct {
		ID uint32
	}

	deleteTaskInteractor struct {
		repo       entities.TaskRepository
		ctxTimeout time.Duration
	}
)

func NewDeleteTaskInteractor(repo entities.TaskRepository, ctxTimeout time.Duration) DeleteTaskUseCase {
	return &deleteTaskInteractor{
		repo:       repo,
		ctxTimeout: ctxTimeout,
	}
}

func (uc deleteTaskInteractor) Execute(ctx context.Context, input DeleteTaskInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	err := uc.repo.DeleteTask(ctx, input.ID)
	if err != nil {
		return err
	}

	return nil
}
