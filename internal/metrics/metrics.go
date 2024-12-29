package metrics

import (
	"net/http"

	"github.com/OwodDEV/crypto-service/internal/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Config *config.Config
}

func NewMetrics(cfg *config.Config) (metrics *Metrics, err error) {
	metrics = &Metrics{
		Config: cfg,
	}
	return
}

func (s *Metrics) PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
