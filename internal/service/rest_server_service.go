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
	CmdExecStatusService *CmdExecStatusService
	URLPath              string
	Port                 string
	Server               *http.Server
	RetSuccess           string
	RetFailed            string
}

func NewRestServerService(css *CmdExecStatusService, urlPath, port, retSuccess, retFailed string) *RestServerService {
	return &RestServerService{CmdExecStatusService: css, URLPath: urlPath, Port: port, RetSuccess: retSuccess, RetFailed: retFailed}
}

func (rss *RestServerService) StartServer() error {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(rss.URLPath,
			func(w rest.ResponseWriter, req *rest.Request) {
				var outputVal string
				switch rss.CmdExecStatusService.CmdExecStatus.Status {
				case domain.Failed:
					outputVal = rss.RetFailed
				case domain.Success:
					outputVal = rss.RetSuccess

				}

				if err := w.WriteJson(&OutputSchema{
					Result: outputVal,
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
