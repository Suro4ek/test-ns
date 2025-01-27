package usecase

import (
	"context"
	"test-ns/internal/entities"
	"time"
)

type (
	CreateTaskUseCase interface {
		Execute(ctx context.Context, input CreateTaskInput) (CreateTaskOutput, error)
	}

	CreateTaskInput struct {
		Title       string
		Description string
		Status      string
	}

	CreateTaskOutput struct {
		ID          uint32    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	CreateTaskPresenter interface {
		Output(task entities.Task) CreateTaskOutput
	}

	createTaskInteractor struct {
		repo       entities.TaskRepository
		presenter  CreateTaskPresenter
		ctxTimeout time.Duration
	}
)

func NewCreateTaskInteractor(repo entities.TaskRepository, presenter CreateTaskPresenter, ctxTimeout time.Duration) CreateTaskUseCase {
	return &createTaskInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: ctxTimeout,
	}
}

func (uc createTaskInteractor) Execute(ctx context.Context, input CreateTaskInput) (CreateTaskOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	task := entities.NewTaskCreate(input.Title, input.Description, input.Status, time.Now())
	task, err := uc.repo.Create(ctx, task)
	if err != nil {
		return CreateTaskOutput{}, err
	}

	return uc.presenter.Output(task), nil
}
