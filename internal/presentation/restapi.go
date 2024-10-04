package presentation

import (
	"context"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// 出力の定義
type OutputSchema struct {
	Result string
}

type RestAPIServer struct {
	URLPath string
	Port    string
	Server  *http.Server
	Handler *Handler
}

func NewRestAPIServer(handler *Handler, urlPath, port string) *RestAPIServer {
	return &RestAPIServer{Handler: handler, URLPath: urlPath, Port: port}
}

func (ras *RestAPIServer) StartServer() error {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(ras.URLPath, ras.Handler.HealthdHandler),
	)
	if err != nil {
		return err
	}
	api.SetApp(router)

	server := &http.Server{
		Addr:    ras.Port,
		Handler: api.MakeHandler(),
	}
	ras.Server = server

	log.Println("HttpServer started")

	return server.ListenAndServe()
}

func (ras *RestAPIServer) ShutdownServer() error {
	return ras.Server.Shutdown(context.Background())
}
