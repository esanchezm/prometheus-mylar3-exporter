// Copyright 2024, Esteban Sanchez

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package exporter

import (
	"context"
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type versionCollector struct {
	ctx  context.Context
	base *baseCollector
}

var (
	metricsUp = prometheus.NewDesc(
		prometheus.BuildFQName(exporterPrefix, "", "up"),
		"Whether the mylar3 server is answering requests from this exporter. "+
			"A `version` label with the server version is added.",
		[]string{"server", "version"},
		nil,
	)
)

func newServerCollector(ctx context.Context, client *mylar3Client, logger *logrus.Logger) *versionCollector {
	return &versionCollector{
		ctx:  ctx,
		base: newBaseCollector(client, logger),
	}
}

func (c *versionCollector) Describe(ch chan<- *prometheus.Desc) {
	c.base.Describe(ch, c.collect)
}

func (c *versionCollector) Collect(ch chan<- prometheus.Metric) {
	c.base.Collect(ch)
}

func (c *versionCollector) collect(ch chan<- prometheus.Metric) {
	client := c.base.client
	logger := c.base.logger

	result := make(chan Mylar3RawResponse)

	go func() {
		response, err := client.CallCommand("getVersion", nil)
		result <- Mylar3RawResponse{Data: response, Err: err}
	}()

	response := <-result
	version := "unknown"
	value := 0.0
	if response.Err == nil {
		logger.Debug("Version info:", string(response.Data))

		var data map[string]interface{}

		err := json.Unmarshal(response.Data, &data)
		if err != nil {
			logger.Errorf("Error unmarshalling version info: %s", err)
			return
		}

		if data["success"].(bool) {
			versionInfo := data["data"].(map[string]interface{})
			version = versionInfo["current_version"].(string)
			value = 1
		}
	}

	ch <- prometheus.MustNewConstMetric(
		metricsUp,
		prometheus.GaugeValue,
		value,
		client.opts.URI,
		version,
	)
}
