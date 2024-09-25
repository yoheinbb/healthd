package util

import (
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
}

func (hs *HttpServer) Start() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(hs.URLPath,
			func(w rest.ResponseWriter, req *rest.Request) {
				if err := w.WriteJson(&OutputSchema{
					hs.Status.Status,
				}); err != nil {
					log.Fatal(err)
				}
			}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/", api.MakeHandler())

	log.Println("HttpServer started")

	log.Fatal(http.ListenAndServe(hs.Port, nil))
}

func NewHttpServer(ss *ServiceStatus, gconfig *GlobalConfig) *HttpServer {

	return &HttpServer{Status: ss, URLPath: gconfig.URLPath, Port: gconfig.Port}
}
