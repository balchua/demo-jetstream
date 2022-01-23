package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	histogramVec *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	jobsDurationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "jobs_duration_seconds",
			Help: "Jobs duration distribution",
		},
		[]string{"job_type"},
	)
	prometheus.MustRegister(jobsDurationHistogram)

	return &Metrics{
		histogramVec: jobsDurationHistogram,
	}
}

func (m *Metrics) Observe(duration time.Duration) {
	m.histogramVec.WithLabelValues("job_duration").Observe(duration.Seconds())
}

func (m *Metrics) Start() {
	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(":2112", nil))
}
