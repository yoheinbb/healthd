package presentation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewAPIServer(logger *slog.Logger, urlPath, port string, handler *Handler) (*http.Server, error) {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(ginLogFormat(logger))
	engine.GET(urlPath, handler.HealthdHandler)
	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return server, nil
}

// Ginのミドルウェア
func ginLogFormat(logger *slog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("gin-request",
			slog.String("time", param.TimeStamp.Format(time.RFC3339)),
			slog.Int("status", param.StatusCode),
			slog.String("latency", param.Latency.String()),
			slog.String("client_ip", param.ClientIP),
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.String("error", param.ErrorMessage),
		)
		return ""
	})
}
