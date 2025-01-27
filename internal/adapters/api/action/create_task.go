package action

import (
	"encoding/json"
	"net/http"
	"test-ns/internal/adapters/api/response"
	"test-ns/internal/usecase"
)

type CreateTaskAction struct {
	uc usecase.CreateTaskUseCase
}

func NewCreateTaskAction(uc usecase.CreateTaskUseCase) CreateTaskAction {
	return CreateTaskAction{uc: uc}
}

func (a CreateTaskAction) Execute(w http.ResponseWriter, r *http.Request) {

	var input usecase.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.NewError(err, http.StatusBadRequest).Send(w)
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		response.NewError(err, http.StatusInternalServerError).Send(w)
	}

	response.NewSuccess(output, http.StatusOK).Send(w)
}
