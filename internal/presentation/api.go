package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gin-gonic/gin"
)

func NewRestAPIServer(urlPath, port string, handler *Handler) (*http.Server, error) {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(urlPath, handler.HealthdHandler),
	)
	if err != nil {
		return nil, err
	}
	api.SetApp(router)

	server := &http.Server{
		Addr:    port,
		Handler: api.MakeHandler(),
	}

	return server, nil
}

func NewAPIServer(urlPath, port string, handler *Handler) (*http.Server, error) {

	engine := gin.Default()
	engine.GET(urlPath, handler.GinHealthdHandler)
	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return server, nil
}
