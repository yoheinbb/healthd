package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
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
