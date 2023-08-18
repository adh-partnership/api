/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package metrics

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/adh-partnership/api/pkg/logger"
)

type MetricType int

const (
	None MetricType = iota
	Counter
	Gauge
	Histogram
	Summary

	defaultMetricPath = "/metrics"
	defaultMetricPort = 8080
	defaultSlowTime   = int32(5)

	MetricRequestTotal    = "gin_request_total"
	MetricRequestUVTotal  = "gin_request_uv_total"
	MetricURIRequestTotal = "gin_uri_request_total"
	MetricRequestBody     = "gin_request_body_total"
	MetricResponseBody    = "gin_response_body_total"
	MetricRequestDuration = "gin_request_duration"
	MetricSlowRequest     = "gin_slow_request_total"
)

var log = logger.Logger.WithField("component", "network")

var (
	defaultDuration = []float64{0.1, 0.3, 1.2, 5, 10}
	monitor         *Monitor

	promTypeHandler = map[MetricType]func(metric *Metric) error{
		Counter:   counterHandler,
		Gauge:     gaugeHandler,
		Histogram: histogramHandler,
		Summary:   summaryHandler,
	}
)

type Monitor struct {
	slowTime    int32
	metricPath  string
	metricPort  int
	reqDuration []float64
	metrics     map[string]*Metric
}

// GetMonitor used to get global Monitor object,
// this function returns a singleton object.
func GetMonitor() *Monitor {
	if monitor == nil {
		monitor = &Monitor{
			metricPath:  defaultMetricPath,
			metricPort:  defaultMetricPort,
			slowTime:    defaultSlowTime,
			reqDuration: defaultDuration,
			metrics:     make(map[string]*Metric),
		}
	}
	return monitor
}

var httpServerOnce sync.Once

func (m *Monitor) Start() {
	go func(m *Monitor) {
		httpServerOnce.Do(func() {
			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(fmt.Sprintf(":%d", m.metricPort), mux)
			if err != nil {
				log.Errorf("Error starting metrics server: %s", err)
				panic(err)
			}
		})
	}(m)
}

// GetMetric used to get metric object by metric_name.
func (m *Monitor) GetMetric(name string) *Metric {
	if metric, ok := m.metrics[name]; ok {
		return metric
	}
	return &Metric{}
}

// SetMetricPath sets the metricPath property. metricPath is used for Prometheus
// to get gin server monitoring data.
func (m *Monitor) SetMetricPath(path string) {
	m.metricPath = path
}

// SetMetricPort sets the metricPort property. metricPort is to expose the metrics
func (m *Monitor) SetMetricPort(port int) {
	m.metricPort = port
}

// SetSlowTime sets the slowTime property. slowTime is used to determine whether
// the request is slow. For "gin_slow_request_total" metric.
func (m *Monitor) SetSlowTime(slowTime int32) {
	m.slowTime = slowTime
}

// SetDuration sets the reqDuration property. reqDuration is used to ginRequestDuration
// metric buckets.
func (m *Monitor) SetDuration(duration []float64) {
	m.reqDuration = duration
}

// AddMetric add custom monitor metric.
func (m *Monitor) AddMetric(metric *Metric) error {
	if _, ok := m.metrics[metric.Name]; ok {
		return errors.Errorf("metric '%s' is existed", metric.Name)
	}

	if metric.Name == "" {
		return errors.Errorf("metric name cannot be empty")
	}
	if f, ok := promTypeHandler[metric.Type]; ok {
		if err := f(metric); err != nil {
			return err
		}

		prometheus.MustRegister(metric.vec)
		m.metrics[metric.Name] = metric
		return nil
	}
	return errors.Errorf("metric type '%d' does not exist", metric.Type)
}

// nolint:unparam // error is used for histogram and summary
func counterHandler(metric *Metric) error {
	metric.vec = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: metric.Name, Help: metric.Description},
		metric.Labels,
	)
	return nil
}

// nolint:unparam // error is used for histogram and summary
func gaugeHandler(metric *Metric) error {
	metric.vec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: metric.Name, Help: metric.Description},
		metric.Labels,
	)
	return nil
}

func histogramHandler(metric *Metric) error {
	if len(metric.Buckets) == 0 {
		return errors.Errorf("metric '%s' is histogram type, cannot lose bucket param", metric.Name)
	}
	metric.vec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: metric.Name, Help: metric.Description, Buckets: metric.Buckets},
		metric.Labels,
	)
	return nil
}

func summaryHandler(metric *Metric) error {
	if len(metric.Objectives) == 0 {
		return errors.Errorf("metric '%s' is summary type, cannot lose objectives param", metric.Name)
	}
	metric.vec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Name: metric.Name, Help: metric.Description, Objectives: metric.Objectives},
		metric.Labels,
	)
	return nil
}
