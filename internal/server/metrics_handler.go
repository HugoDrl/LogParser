package server

import (
	"encoding/json"
	"net/http"
)

func (d *DataLayer) getMetrics(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(d.Metrics)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (d *DataLayer) AttachMetricsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /metrics", d.getMetrics)
}
