package presenter

import (
	"test-ns/internal/entities"
	"test-ns/internal/usecase"
)

type CreateTaskPresenter struct{}

func NewCreateTaskPresenter() usecase.CreateTaskPresenter {
	return CreateTaskPresenter{}
}

func (p CreateTaskPresenter) Output(task entities.Task) usecase.CreateTaskOutput {
	return usecase.CreateTaskOutput{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		CreatedAt:   task.CreatedAt(),
	}
}
