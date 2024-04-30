package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type baseCollector struct {
	client       *mylar3Client
	logger       *logrus.Logger
	metricsCache []prometheus.Metric
}

func newBaseCollector(client *mylar3Client, logger *logrus.Logger) *baseCollector {
	return &baseCollector{
		client: client,
		logger: logger,
	}
}

func (d *baseCollector) Describe(ch chan<- *prometheus.Desc, collect func(mCh chan<- prometheus.Metric)) {
	d.metricsCache = make([]prometheus.Metric, 0, 1000)

	metrics := make(chan prometheus.Metric)
	go func() {
		collect(metrics)
		close(metrics)
	}()

	for m := range metrics {
		d.metricsCache = append(d.metricsCache, m)
		ch <- m.Desc()
	}
}

func (d *baseCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range d.metricsCache {
		ch <- metric
	}
}
