package usecase

import (
	"context"
	"test-ns/internal/adapters/repo"
	"test-ns/internal/entities"
	"time"
)

type (
	UpdateTaskUseCase interface {
		Execute(ctx context.Context, input UpdateTaskInput) (UpdateTaskOutput, error)
	}

	UpdateTaskInput struct {
		ID          uint32
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	UpdateTaskOutput struct {
		ID          uint32     `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}

	UpdateTaskPresenter interface {
		Output(task entities.Task) UpdateTaskOutput
	}

	updateTaskInteractor struct {
		transaction repo.GTransaction
		repo        entities.TaskRepository
		presenter   UpdateTaskPresenter
		ctxTimeout  time.Duration
	}
)

func NewUpdateTaskInteractor(transaction repo.GTransaction, repo entities.TaskRepository, presenter UpdateTaskPresenter, ctxTimeout time.Duration) UpdateTaskUseCase {
	return &updateTaskInteractor{
		transaction: transaction,
		repo:        repo,
		presenter:   presenter,
		ctxTimeout:  ctxTimeout,
	}
}

func (uc updateTaskInteractor) Execute(ctx context.Context, input UpdateTaskInput) (UpdateTaskOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	ctx = uc.transaction.Begin(ctx)

	task, err := uc.repo.GetTask(ctx, input.ID)
	if err != nil {
		uc.transaction.Rollback(ctx)
		return UpdateTaskOutput{}, err
	}

	time := time.Now()
	task.SetUpdatedAt(&time)
	task.SetTitle(input.Title)
	task.SetDescription(input.Description)
	err = task.SetStatus(input.Status)
	if err != nil {
		uc.transaction.Rollback(ctx)
		return UpdateTaskOutput{}, err
	}
	task, err = uc.repo.Update(ctx, task)
	if err != nil {
		uc.transaction.Rollback(ctx)
		return UpdateTaskOutput{}, err
	}

	uc.transaction.Commit(ctx)
	return uc.presenter.Output(task), nil
}
