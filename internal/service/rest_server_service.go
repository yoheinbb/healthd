package service

import (
	"context"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yoheinbb/healthd/internal/domain"
)

// 出力の定義
type OutputSchema struct {
	Result string
}

type RestServerService struct {
	Status     *domain.Status
	URLPath    string
	Port       string
	Server     *http.Server
	RetSuccess string
	RetFailed  string
}

func NewRestServerService(status *domain.Status, urlPath, port, retSuccess, retFailed string) *RestServerService {
	return &RestServerService{Status: status, URLPath: urlPath, Port: port, RetSuccess: retSuccess, RetFailed: retFailed}
}

func (rss *RestServerService) StartServer() error {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(rss.URLPath, rss.healthdHandler),
	)
	if err != nil {
		return err
	}
	api.SetApp(router)

	server := &http.Server{
		Addr:    rss.Port,
		Handler: api.MakeHandler(),
	}
	rss.Server = server

	log.Println("HttpServer started")

	return server.ListenAndServe()
}

func (rss *RestServerService) ShutdownServer() error {
	return rss.Server.Shutdown(context.Background())
}

func (rss *RestServerService) healthdHandler(w rest.ResponseWriter, _ *rest.Request) {
	var outputVal string
	switch rss.Status.GetStatus() {
	case domain.Success:
		outputVal = rss.RetSuccess
	case domain.Failed:
		outputVal = rss.RetFailed
	default:
		outputVal = rss.RetFailed
	}

	if err := w.WriteJson(&OutputSchema{
		Result: outputVal,
	}); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
