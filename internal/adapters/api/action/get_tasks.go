package action

import (
	"net/http"
	"test-ns/internal/adapters/api/response"
	"test-ns/internal/usecase"
)

type GetTasksAction struct {
	uc usecase.GetTasksUseCase
}

func NewGetTasksAction(uc usecase.GetTasksUseCase) GetTasksAction {
	return GetTasksAction{uc: uc}
}

func (a GetTasksAction) Execute(w http.ResponseWriter, r *http.Request) {

	var input usecase.GetTasksInput

	var keysStr = r.URL.Query()
	sort, ok := keysStr["sort"]
	if ok {
		sortStr := sort[0]
		input.Sort = &sortStr
	}

	status, ok := keysStr["status"]
	if ok {
		statusStr := status[0]
		input.Status = &statusStr
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		response.NewError(err, http.StatusInternalServerError).Send(w)
	}

	response.NewSuccess(output, http.StatusCreated).Send(w)
}
