package server

import (
	"encoding/json"
	"net/http"

	"github.com/HugoDrl/zebra/internal/analyser"
)

func (d *DataLayer) calculateMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := analyser.AnalyseLogs(d.Logs, d.Errs, d.AnalyseSettings)

	payload, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (d *DataLayer) AttachMetricsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /metrics", d.calculateMetrics)
}
