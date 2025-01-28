package action

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-ns/internal/adapters/api/response"
	"test-ns/internal/usecase"
)

type UpdateTaskAction struct {
	uc usecase.UpdateTaskUseCase
}

func NewUpdateTaskAction(uc usecase.UpdateTaskUseCase) UpdateTaskAction {
	return UpdateTaskAction{uc: uc}
}

func (a UpdateTaskAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.UpdateTaskInput

	var keysStr = r.URL.Query()
	id, err := strconv.Atoi(keysStr["id"][0])
	if err != nil {
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	input.ID = uint32(id)

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	response.NewSuccess(output, http.StatusCreated).Send(w)
}
