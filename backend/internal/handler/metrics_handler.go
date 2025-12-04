package handler

import (
	"net/http"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsHandler handles metrics endpoints.
type MetricsHandler struct {
	startTime     time.Time
	requestCount  atomic.Uint64
	errorCount    atomic.Uint64
}

// MetricsResponse represents the metrics response.
type MetricsResponse struct {
	Uptime        string        `json:"uptime"`
	UptimeSeconds float64       `json:"uptime_seconds"`
	Requests      uint64        `json:"requests_total"`
	Errors        uint64        `json:"errors_total"`
	Memory        MemoryMetrics `json:"memory"`
	Goroutines    int           `json:"goroutines"`
}

// MemoryMetrics contains memory-related metrics.
type MemoryMetrics struct {
	Alloc      uint64 `json:"alloc_bytes"`
	TotalAlloc uint64 `json:"total_alloc_bytes"`
	Sys        uint64 `json:"sys_bytes"`
	NumGC      uint32 `json:"num_gc"`
}

// NewMetricsHandler creates a new MetricsHandler instance.
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		startTime: time.Now(),
	}
}

// IncrementRequests increments the request counter.
func (h *MetricsHandler) IncrementRequests() {
	h.requestCount.Add(1)
}

// IncrementErrors increments the error counter.
func (h *MetricsHandler) IncrementErrors() {
	h.errorCount.Add(1)
}

// Metrics returns application metrics.
// @Summary Get application metrics
// @Description Returns application metrics including uptime, memory usage, and request counts
// @Tags monitoring
// @Produce json
// @Success 200 {object} MetricsResponse
// @Router /metrics [get]
func (h *MetricsHandler) Metrics(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(h.startTime)

	c.JSON(http.StatusOK, MetricsResponse{
		Uptime:        uptime.String(),
		UptimeSeconds: uptime.Seconds(),
		Requests:      h.requestCount.Load(),
		Errors:        h.errorCount.Load(),
		Memory: MemoryMetrics{
			Alloc:      memStats.Alloc,
			TotalAlloc: memStats.TotalAlloc,
			Sys:        memStats.Sys,
			NumGC:      memStats.NumGC,
		},
		Goroutines: runtime.NumGoroutine(),
	})
}

// PrometheusMetrics returns metrics in Prometheus format.
// @Summary Get Prometheus metrics
// @Description Returns metrics in Prometheus text format
// @Tags monitoring
// @Produce text/plain
// @Success 200 {string} string
// @Router /metrics/prometheus [get]
func (h *MetricsHandler) PrometheusMetrics(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(h.startTime)

	// Simple Prometheus-compatible text format
	metrics := ""
	metrics += "# HELP superdash_uptime_seconds Time since service started\n"
	metrics += "# TYPE superdash_uptime_seconds gauge\n"
	metrics += "superdash_uptime_seconds " + formatFloat(uptime.Seconds()) + "\n"
	metrics += "\n"
	metrics += "# HELP superdash_requests_total Total number of requests\n"
	metrics += "# TYPE superdash_requests_total counter\n"
	metrics += "superdash_requests_total " + formatUint64(h.requestCount.Load()) + "\n"
	metrics += "\n"
	metrics += "# HELP superdash_errors_total Total number of errors\n"
	metrics += "# TYPE superdash_errors_total counter\n"
	metrics += "superdash_errors_total " + formatUint64(h.errorCount.Load()) + "\n"
	metrics += "\n"
	metrics += "# HELP superdash_memory_alloc_bytes Current memory allocation\n"
	metrics += "# TYPE superdash_memory_alloc_bytes gauge\n"
	metrics += "superdash_memory_alloc_bytes " + formatUint64(memStats.Alloc) + "\n"
	metrics += "\n"
	metrics += "# HELP superdash_goroutines Current number of goroutines\n"
	metrics += "# TYPE superdash_goroutines gauge\n"
	metrics += "superdash_goroutines " + formatInt(runtime.NumGoroutine()) + "\n"

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, metrics)
}

// RegisterMetricsRoutes registers metrics routes.
func (h *MetricsHandler) RegisterMetricsRoutes(r *gin.Engine) {
	r.GET("/metrics", h.Metrics)
	r.GET("/metrics/prometheus", h.PrometheusMetrics)
}

// MetricsMiddleware returns a middleware that tracks requests.
func (h *MetricsHandler) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.IncrementRequests()
		c.Next()
		if c.Writer.Status() >= 400 {
			h.IncrementErrors()
		}
	}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func formatUint64(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}
