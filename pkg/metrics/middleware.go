package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/bloom"
)

var bloomFilter *bloom.BloomFilter

// Use set gin metrics middleware
func (m *Monitor) Use(r gin.IRoutes) {
	m.initGinMetrics()

	r.Use(m.monitorInterceptor)
}

// initGinMetrics used to init gin metrics
func (m *Monitor) initGinMetrics() {
	bloomFilter = bloom.NewBloomFilter()

	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricRequestTotal,
		Description: "total number of requests received",
		Labels:      nil,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricRequestUVTotal,
		Description: "total number of unique requests received",
		Labels:      nil,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricURIRequestTotal,
		Description: "total number of requests received by uri",
		Labels:      []string{"uri", "method", "code"},
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricRequestBody,
		Description: "total received request body size, unit byte",
		Labels:      nil,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricResponseBody,
		Description: "total sent response body size, unit byte",
		Labels:      nil,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Histogram,
		Name:        MetricRequestDuration,
		Description: "the time server took to handle the request",
		Labels:      []string{"uri"},
		Buckets:     m.reqDuration,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        MetricSlowRequest,
		Description: fmt.Sprintf("slow request counter, t=%d", m.slowTime),
		Labels:      []string{"uri", "method", "code"},
	})
}

// monitorInterceptor as gin monitor middleware.
func (m *Monitor) monitorInterceptor(ctx *gin.Context) {
	startTime := time.Now()

	// execute normal process.
	ctx.Next()

	// after request
	m.ginMetricHandle(ctx, startTime)
}

func (m *Monitor) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer
	latency := time.Since(start)

	log.Tracef("Handling metric for %s", ctx.FullPath())

	// set request total
	_ = m.GetMetric(MetricRequestTotal).Inc(nil)

	// set uv
	if clientIP := ctx.ClientIP(); !bloomFilter.Contains(clientIP) {
		bloomFilter.Add(clientIP)
		_ = m.GetMetric(MetricRequestUVTotal).Inc(nil)
	}

	// set uri request total
	_ = m.GetMetric(MetricURIRequestTotal).Inc([]string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status())})

	// set request body size
	// since r.ContentLength can be negative (in some occasions) guard the operation
	if r.ContentLength >= 0 {
		_ = m.GetMetric(MetricRequestBody).Add(nil, float64(r.ContentLength))
	}

	// set slow request
	if int32(latency.Seconds()) > m.slowTime {
		_ = m.GetMetric(MetricSlowRequest).Inc([]string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status())})
	}

	// set request duration
	_ = m.GetMetric(MetricRequestDuration).Observe([]string{ctx.FullPath()}, latency.Seconds())

	// set response size
	if w.Size() > 0 {
		_ = m.GetMetric(MetricResponseBody).Add(nil, float64(w.Size()))
	}
}
