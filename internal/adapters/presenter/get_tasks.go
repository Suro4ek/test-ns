package presenter

import (
	"test-ns/internal/entities"
	"test-ns/internal/usecase"
)

type GetTasksPresenter struct{}

func NewGetTasksPresenter() usecase.GetTasksPresenter {
	return GetTasksPresenter{}
}

func (p GetTasksPresenter) Output(tasks []entities.Task) usecase.GetTasksListOutput {
	var tasksOutput []usecase.GetTasksOutput
	for _, task := range tasks {
		tasksOutput = append(tasksOutput, usecase.GetTasksOutput{
			ID:          task.ID(),
			Title:       task.Title(),
			Description: task.Description(),
			Status:      task.Status(),
			CreatedAt:   task.CreatedAt(),
			UpdatedAt:   task.UpdatedAt(),
		})
	}

	return usecase.GetTasksListOutput{
		Tasks: tasksOutput,
	}
}
