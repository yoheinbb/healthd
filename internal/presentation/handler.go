package presentation

import (
	"github.com/gin-gonic/gin"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Handler struct {
	status     usecase.IStatus
	retSuccess string
	retFailed  string
}

func NewHandler(status usecase.IStatus, retSuccess, retFailed string) *Handler {
	return &Handler{status: status, retSuccess: retSuccess, retFailed: retFailed}
}

type OutputSchema struct {
	Result string
}

func (h *Handler) HealthdHandler(c *gin.Context) {
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
