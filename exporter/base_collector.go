//

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
