package asr_gateway

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mDecisions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "asrgw_decisions_total",
		Help: "Total decisions emitted",
	}, []string{"result"})
	mLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "asrgw_decision_latency_ms",
		Help:    "Decision latency in ms",
		Buckets: []float64{100, 200, 300, 500, 800, 1200, 2000, 5000},
	})
	mActiveSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "asrgw_active_sessions",
		Help: "Active WS sessions from mod_audio_fork",
	})
	mForwardedFrames = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "asrgw_forwarded_pcm_frames_total",
		Help: "Total PCM frames forwarded to FunASR",
	})
)

func init() {
	prometheus.MustRegister(mDecisions, mLatency, mActiveSessions, mForwardedFrames)
}

func (s *Server) MetricsHandler() http.Handler { return promhttp.Handler() }