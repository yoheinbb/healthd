package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Handler struct {
	Status  *usecase.Status
	gConfig *util.GlobalConfig
}

func NewHandler(status *usecase.Status, gConfig *util.GlobalConfig) *Handler {
	return &Handler{Status: status, gConfig: gConfig}
}

func (h *Handler) HealthdHandler(w rest.ResponseWriter, _ *rest.Request) {
	var outputVal string
	switch h.Status.GetStatus() {
	case constant.Success:
		outputVal = h.gConfig.RetSuccess
	case constant.Failed:
		outputVal = h.gConfig.RetFailed
	default:
		outputVal = h.gConfig.RetFailed
	}

	if err := w.WriteJson(&OutputSchema{
		Result: outputVal,
	}); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
