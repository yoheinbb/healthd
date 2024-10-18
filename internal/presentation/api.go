package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAPIServer(urlPath, port string, handler *Handler) (*http.Server, error) {

	engine := gin.Default()
	engine.GET(urlPath, handler.HealthdHandler)
	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return server, nil
}
