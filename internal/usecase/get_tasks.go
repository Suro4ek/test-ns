package usecase

import (
	"context"
	"test-ns/internal/entities"
	"time"
)

type (
	GetTasksUseCase interface {
		Execute(ctx context.Context, input GetTasksInput) (GetTasksListOutput, error)
	}

	GetTasksInput struct {
		Status *string
		Sort   *string
	}

	GetTasksOutput struct {
		ID          uint32     `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}

	GetTasksListOutput struct {
		Tasks []GetTasksOutput `json:"tasks"`
	}

	GetTasksPresenter interface {
		Output(tasks []entities.Task) GetTasksListOutput
	}

	getTasksInteractor struct {
		repo       entities.TaskRepository
		presenter  GetTasksPresenter
		ctxTimeout time.Duration
	}
)

func NewTaskInteractor(repo entities.TaskRepository, presenter GetTasksPresenter, ctxTimeout time.Duration) GetTasksUseCase {
	return &getTasksInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: ctxTimeout,
	}
}

func (uc getTasksInteractor) Execute(ctx context.Context, input GetTasksInput) (GetTasksListOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	if input.Status != nil {
		status := *input.Status
		if status != "todo" && status != "in_progress" && status != "done" {
			return GetTasksListOutput{}, entities.ErrorStatus
		}
	}

	if input.Sort != nil {
		sort := *input.Sort
		if sort != "title" && sort != "description" && sort != "status" && sort != "created_at" {
			return GetTasksListOutput{}, entities.ErrorSort
		}
	}

	tasks, err := uc.repo.GetTasks(ctx, input.Status, input.Sort)
	if err != nil {
		return GetTasksListOutput{}, err
	}

	return uc.presenter.Output(tasks), nil

}
