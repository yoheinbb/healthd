package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gin-gonic/gin"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Handler struct {
	status     *usecase.Status
	retSuccess string
	retFailed  string
}

func NewHandler(status *usecase.Status, retSuccess, retFailed string) *Handler {
	return &Handler{status: status, retSuccess: retSuccess, retFailed: retFailed}
}

type OutputSchema struct {
	Result string
}

func (h *Handler) HealthdHandler(w rest.ResponseWriter, _ *rest.Request) {
	var outputVal string
	switch h.status.GetStatus() {
	case constant.SUCCESS:
		outputVal = h.retSuccess
	case constant.FAILED:
		outputVal = h.retFailed
	default:
		outputVal = h.retFailed
	}

	if err := w.WriteJson(&OutputSchema{
		Result: outputVal,
	}); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GinHealthdHandler(c *gin.Context) {
	var outputVal string
	switch h.status.GetStatus() {
	case constant.SUCCESS:
		outputVal = h.retSuccess
	case constant.FAILED:
		outputVal = h.retFailed
	default:
		outputVal = h.retFailed
	}

	c.IndentedJSON(200, &OutputSchema{
		Result: outputVal,
	})
}
