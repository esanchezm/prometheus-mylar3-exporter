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

type wantedCollector struct {
	ctx  context.Context
	base *baseCollector
}

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

	result := make(chan Mylar3RawResponse)

	go func() {
		response, err := client.CallCommand("getWanted", nil)
		result <- Mylar3RawResponse{Data: response, Err: err}
	}()

	response := <-result

	if response.Err == nil {
		logger.Debug("getWanted response:", string(response.Data))
		err := json.Unmarshal(response.Data, &data)
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
					Name: prometheus.BuildFQName(exporterPrefix, "", "wanted_count"),
					Help: "Number of issues in `Wanted` state. Labels `publisher`, `name` and `year` are added",
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
