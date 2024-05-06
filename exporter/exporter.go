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

func New(opts *Mylar3Opts, logger *logrus.Logger) *Exporter {
	return &Exporter{
		client: newMylar3Client(opts, logger),
		logger: logger,
		opts:   opts,
	}
}

func (e *Exporter) makeRegistry(ctx context.Context) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(newServerCollector(ctx, e.client, e.logger))
	registry.MustRegister(newIndexCollector(ctx, e.client, e.logger))
	registry.MustRegister(newWantedCollector(ctx, e.client, e.logger))

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
