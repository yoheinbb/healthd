package util

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
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

func (self *HttpServer) Start() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(self.URLPath,
			func(w rest.ResponseWriter, req *rest.Request) {
				w.WriteJson(&OutputSchema{
					self.Status.Status,
				})
			}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/", api.MakeHandler())

	log.Println("HttpServer started")

	log.Fatal(http.ListenAndServe(self.Port, nil))
}

func NewHttpServer(ss *ServiceStatus, gconfig *GlobalConfig) *HttpServer {

	return &HttpServer{Status: ss, URLPath: gconfig.URLPath, Port: gconfig.Port}
}
