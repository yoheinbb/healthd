package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Handler struct {
	Status     *usecase.Status
	RetSuccess string
	RetFailed  string
}

func NewHandler(status *usecase.Status, retSuccess, retFailed string) *Handler {
	return &Handler{Status: status, RetSuccess: retSuccess, RetFailed: retFailed}
}

func (h *Handler) HealthdHandler(w rest.ResponseWriter, _ *rest.Request) {
	var outputVal string
	switch h.Status.GetStatus() {
	case constant.Success:
		outputVal = h.RetSuccess
	case constant.Failed:
		outputVal = h.RetFailed
	default:
		outputVal = h.RetFailed
	}

	if err := w.WriteJson(&OutputSchema{
		Result: outputVal,
	}); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
