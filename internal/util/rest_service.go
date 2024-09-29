package util

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

type HttpServer struct {
	Status  *ServiceStatus
	URLPath string
	Port    string
	Server  *http.Server
}

func NewHttpServer(ss *ServiceStatus, urlPath, port string) *HttpServer {

	return &HttpServer{Status: ss, URLPath: urlPath, Port: port}
}

func (hs *HttpServer) Start() error {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(hs.URLPath,
			func(w rest.ResponseWriter, req *rest.Request) {
				if err := w.WriteJson(&OutputSchema{
					hs.Status.Status,
				}); err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}),
	)
	if err != nil {
		return err
	}
	api.SetApp(router)

	server := &http.Server{
		Addr:    hs.Port,
		Handler: api.MakeHandler(),
	}
	hs.Server = server

	log.Println("HttpServer started")

	return server.ListenAndServe()
}

func (hs *HttpServer) Shutdown() error {
	return hs.Server.Shutdown(context.Background())
}
