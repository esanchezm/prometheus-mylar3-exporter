package exporter

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Exporter holds Exporter methods and attributes.
type Exporter struct {
	client *mylar3Client
	logger *logrus.Logger
	opts   *Mylar3Opts
}

func New(opts *Mylar3Opts) *Exporter {
	logger := logrus.New()
	return &Exporter{
		client: newMylar3Client(opts, logger),
		logger: logger,
		opts:   opts,
	}
}

func (e *Exporter) makeRegistry(ctx context.Context) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	serverCollector := newServerCollector(ctx, e.client, e.logger)
	registry.MustRegister(serverCollector)

	return registry
}

func (e *Exporter) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seconds := e.opts.Timeout

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(seconds)*time.Second)
		defer cancel()

		registry := e.makeRegistry(ctx)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
			ErrorLog:      e.logger,
		})

		h.ServeHTTP(w, r)
	})
}
