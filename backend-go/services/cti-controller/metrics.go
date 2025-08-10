package cticontroller

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mAsrCallbacks = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cti_asr_callbacks_total",
		Help: "ASR decisions received",
	}, []string{"result"})
)

func init() { prometheus.MustRegister(mAsrCallbacks) }

func (s *Server) MetricsHandler() http.Handler { return promhttp.Handler() }