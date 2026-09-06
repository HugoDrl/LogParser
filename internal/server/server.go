package server

import (
	"net/http"

	"github.com/HugoDrl/zebra/internal/analyser"
)

type DataLayer struct {
	Metrics analyser.CollectionMetric
}

type HttpServer struct {
	server    *http.Server
	dataLayer *DataLayer
}

func NewServer(data *DataLayer) *HttpServer {
	handler := http.NewServeMux()
	data.AttachMetricsHandler(handler)

	s := &HttpServer{
		server:    &http.Server{Addr: ":8000"},
		dataLayer: data,
	}
	s.server.Handler = handler
	return s
}

func (s *HttpServer) StartServer() error {
	return s.server.ListenAndServe()
}
