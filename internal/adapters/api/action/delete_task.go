package action

import (
	"net/http"
	"strconv"
	"test-ns/internal/adapters/api/response"
	"test-ns/internal/usecase"
)

type DeleteTaskAction struct {
	uc usecase.DeleteTaskUseCase
}

func NewDeleteTaskAction(uc usecase.DeleteTaskUseCase) DeleteTaskAction {
	return DeleteTaskAction{uc: uc}
}

func (a DeleteTaskAction) Execute(w http.ResponseWriter, r *http.Request) {

	var input usecase.DeleteTaskInput
	var keysStr = r.URL.Query()
	id, err := strconv.Atoi(keysStr["id"][0])
	if err != nil {
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	input.ID = uint32(id)

	err = a.uc.Execute(r.Context(), input)
	if err != nil {
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	response.NewSuccess(nil, http.StatusCreated).Send(w)
}
