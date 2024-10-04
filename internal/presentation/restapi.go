package presentation

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yoheinbb/healthd/internal/util"
)

// 出力の定義
type OutputSchema struct {
	Result string
}

func NewRestAPIServer(handler *Handler, gConfig *util.GlobalConfig) (*http.Server, error) {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(gConfig.URLPath, handler.HealthdHandler),
	)
	if err != nil {
		return nil, err
	}
	api.SetApp(router)

	server := &http.Server{
		Addr:    gConfig.Port,
		Handler: api.MakeHandler(),
	}

	return server, nil
}
