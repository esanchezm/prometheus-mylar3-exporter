package exporter

import (
	"context"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type indexCollector struct {
	ctx  context.Context
	base *baseCollector
}

var (
	metricsSeriesCount = prometheus.NewDesc(
		prometheus.BuildFQName(exporterPrefix, "", "series_count"),
		"Number of watched series.",
		[]string{"server", "publisher"},
		nil,
	)

	metricsIssuesCount = prometheus.NewDesc(
		prometheus.BuildFQName(exporterPrefix, "", "issues_count"),
		"Total number of issues available.",
		[]string{"server", "publisher", "name", "year", "status"},
		nil,
	)
)

func newIndexCollector(ctx context.Context, client *mylar3Client, logger *logrus.Logger) *indexCollector {
	return &indexCollector{
		ctx:  ctx,
		base: newBaseCollector(client, logger),
	}
}

func (c *indexCollector) Describe(ch chan<- *prometheus.Desc) {
	c.base.Describe(ch, c.collect)
}

func (c *indexCollector) Collect(ch chan<- prometheus.Metric) {
	c.base.Collect(ch)
}

func (c *indexCollector) collect(ch chan<- prometheus.Metric) {
	var data map[string]interface{}

	client := c.base.client
	logger := c.base.logger

	response, err := client.CallCommand("getIndex", nil)
	if err == nil {
		logger.Debug("getIndex response:", string(response))
		err := json.Unmarshal(response, &data)
		if err != nil {
			logger.Errorf("Error unmarshalling response: %s", err)
			return
		}
	}

	var publisherCounter = make(map[string]float64)
	if data["success"].(bool) {
		for _, v := range data["data"].([]interface{}) {
			publisher := v.(map[string]interface{})["publisher"].(string)
			publisherCounter[publisher]++
		}
	}

	for publisher, value := range publisherCounter {
		ch <- prometheus.MustNewConstMetric(
			metricsSeriesCount,
			prometheus.GaugeValue,
			value,
			client.opts.URI,
			publisher,
		)
	}

	if data["success"].(bool) {
		for _, v := range data["data"].([]interface{}) {
			ch <- prometheus.MustNewConstMetric(
				metricsIssuesCount,
				prometheus.GaugeValue,
				v.(map[string]interface{})["totalIssues"].(float64),
				client.opts.URI,
				v.(map[string]interface{})["publisher"].(string),
				v.(map[string]interface{})["name"].(string),
				v.(map[string]interface{})["year"].(string),
				v.(map[string]interface{})["status"].(string),
			)
		}
	}

}
