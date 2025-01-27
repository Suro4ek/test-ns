package presenter

import (
	"test-ns/internal/entities"
	"test-ns/internal/usecase"
)

type UpdateTaskPresenter struct{}

func NewUpdateTaskPresenter() usecase.UpdateTaskPresenter {
	return UpdateTaskPresenter{}
}

func (p UpdateTaskPresenter) Output(task entities.Task) usecase.UpdateTaskOutput {
	return usecase.UpdateTaskOutput{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}
}
