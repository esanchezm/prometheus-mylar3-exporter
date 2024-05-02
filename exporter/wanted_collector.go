package exporter

import (
	"context"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type wantedCollector struct {
	ctx  context.Context
	base *baseCollector
}

var (
	metricsWantedCount = prometheus.NewDesc(
		prometheus.BuildFQName(exporterPrefix, "", "wanted_count"),
		"Number of issues in `Wanted` state. Labels `publisher`, `name` and `year` are added",
		[]string{"server", "publisher", "name", "year"},
		nil,
	)
)

func newWantedCollector(ctx context.Context, client *mylar3Client, logger *logrus.Logger) *wantedCollector {
	return &wantedCollector{
		ctx:  ctx,
		base: newBaseCollector(client, logger),
	}
}

func (c *wantedCollector) Describe(ch chan<- *prometheus.Desc) {
	c.base.Describe(ch, c.collect)
}

func (c *wantedCollector) Collect(ch chan<- prometheus.Metric) {
	c.base.Collect(ch)
}

func (c *wantedCollector) collect(ch chan<- prometheus.Metric) {
	var data map[string]interface{}

	client := c.base.client
	logger := c.base.logger

	response, err := client.CallCommand("getWanted", nil)
	if err == nil {
		logger.Debug("getWanted response:", string(response))
		err := json.Unmarshal(response, &data)
		if err != nil {
			logger.Errorf("Error unmarshalling response: %s", err)
			return
		}
	}

	metrics := make(map[string]prometheus.Gauge)
	if data["issues"] != nil {
		for _, v := range data["issues"].([]interface{}) {
			// Make string key for publisher, name and year
			publisher := v.(map[string]interface{})["ComicPublisher"].(string)
			name := v.(map[string]interface{})["ComicName"].(string)
			year := v.(map[string]interface{})["ComicYear"].(string)
			key := publisher + name + year
			if metrics[key] == nil {
				metrics[key] = prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: exporterPrefix,
					Name:      prometheus.BuildFQName(exporterPrefix, "", "wanted_count"),
					Help:      "Number of issues in `Wanted` state. Labels `publisher`, `name` and `year` are added",
					ConstLabels: map[string]string{
						"server":    client.opts.URI,
						"publisher": publisher,
						"name":      name,
						"year":      year,
					},
				})
			}
			metrics[key].Inc()
		}
	}

	for _, metrics := range metrics {
		ch <- metrics
	}
}