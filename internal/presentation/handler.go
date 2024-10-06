package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Handler struct {
	status  *usecase.Status
	gConfig *util.GlobalConfig
}

func NewHandler(status *usecase.Status, gConfig *util.GlobalConfig) *Handler {
	return &Handler{status: status, gConfig: gConfig}
}

func (h *Handler) HealthdHandler(w rest.ResponseWriter, _ *rest.Request) {
	var outputVal string
	switch h.status.GetStatus() {
	case constant.SUCCESS:
		outputVal = h.gConfig.RetSuccess
	case constant.FAILED:
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
