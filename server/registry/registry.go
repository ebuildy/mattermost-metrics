package registry

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Registry struct {
	Registry *prometheus.Registry
}


func NewRegistry() *Registry {
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &Registry{
		Registry: reg,
	}
}

func (reg *Registry) HandleHTTP() http.Handler {
	return promhttp.HandlerFor(reg.Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry: reg.Registry,
	})
}