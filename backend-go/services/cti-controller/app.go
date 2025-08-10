package cticontroller

import (
	"log"
	"net/url"
	"os"
	"time"

	cti "github.com/cati/system/internal/cti"
	"github.com/prometheus/client_golang/prometheus"
)

var progressMediaTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "cti_progress_media_total",
	Help: "Total CHANNEL_PROGRESS_MEDIA events",
})

func init() { prometheus.MustRegister(progressMediaTotal) }

type Server struct {
	esl   *cti.Client
	asrGW string // ws base url for audio_fork
	store *Store
}

func NewServer() *Server {
	s := &Server{}
	s.asrGW = os.Getenv("ASR_GATEWAY_WS") // e.g., wss://asr-gateway:10000/stream
	st, err := NewStoreFromEnv()
	if err != nil { log.Printf("pg init error: %v", err) } else { s.store = st }
	return s
}

func (s *Server) StartESL() {
	host := os.Getenv("ESL_HOST")
	port := os.Getenv("ESL_PORT")
	pw := os.Getenv("ESL_PASSWORD")
	if host == "" || port == "" || pw == "" {
		log.Printf("ESL not configured, skip connect")
		return
	}
	addr := host + ":" + port
	h := func(ev cti.Event) {
		if ev.Headers["Content-Type"] == "text/event-plain" {
			if ev.Headers["Event-Name"] == "CHANNEL_PROGRESS_MEDIA" {
				progressMediaTotal.Inc()
				uuid := ev.Headers["Unique-ID"]
				s.handleProgressMedia(uuid)
			}
		}
	}
	c := cti.NewClient(addr, pw, h)
	for {
		if err := c.Connect(); err != nil {
			log.Printf("ESL connect error: %v, retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Printf("ESL connected to %s", addr)
		s.esl = c
		return
	}
}

func (s *Server) handleProgressMedia(uuid string) {
	if s.esl == nil || s.asrGW == "" { return }
	u, _ := url.Parse(s.asrGW)
	q := u.Query()
	q.Set("uuid", uuid)
	u.RawQuery = q.Encode()
	_ = s.esl.UUIDAudioForkStart(uuid, u.String(), "{channels=1,sampling=8000,ptime=20,stream-type=sender}")
}

func (s *Server) KillByDecision(uuid string) {
	if s.esl == nil { return }
	_ = s.esl.UUIDKill(uuid, "NORMAL_CLEARING")
}