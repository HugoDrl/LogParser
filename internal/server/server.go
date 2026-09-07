package server

import (
	"net/http"
	"time"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/parser"
)

type DataLayer struct {
	Logs            []*parser.Log
	Errs            []error
	AnalyseSettings *analyser.AnalyserSettings
}

type HttpServer struct {
	server    *http.Server
	dataLayer *DataLayer
}

func NewServer(data *DataLayer) *HttpServer {
	handler := http.NewServeMux()
	data.AttachMetricsHandler(handler)

	s := &HttpServer{
		server:    &http.Server{Addr: ":8000", ReadHeaderTimeout: 100 * time.Millisecond},
		dataLayer: data,
	}
	s.server.Handler = handler
	return s
}

func (s *HttpServer) StartServer() error {
	return s.server.ListenAndServe()
}
