package asr_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cati/system/internal/asr"
	fa "github.com/cati/system/internal/funasr"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Server struct { phrases atomic.Value /* *asr.Phrases */ }

type Decision struct {
	UUID       string  `json:"uuid"`
	Result     string  `json:"result"`
	Confidence float64 `json:"confidence"`
	LatencyMs  int     `json:"latency_ms"`
	Transcript string  `json:"transcript"`
	Mode       string  `json:"mode"`
	ProofURI   string  `json:"audio_proof_uri"`
	Fallback   bool    `json:"fallback"`
}

func NewServer() (*Server, error) {
	s := &Server{}
	if err := s.loadPhrases(); err != nil { return nil, err }
	return s, nil
}

func (s *Server) loadPhrases() error {
	path := os.Getenv("PHRASES_PATH")
	if path == "" { path = filepath.Join("/", "config", "phrases.yml") }
	phr, err := asr.LoadPhrases(path)
	if err != nil { return fmt.Errorf("load phrases: %w", err) }
	s.phrases.Store(phr)
	return nil
}

func (s *Server) HandleReload(w http.ResponseWriter, r *http.Request) {
	if err := s.loadPhrases(); err != nil { http.Error(w, err.Error(), 500); return }
	w.WriteHeader(200); w.Write([]byte("reloaded"))
}

func (s *Server) callbackCTI(dec Decision) {
	url := os.Getenv("CTI_DECISION_URL"); if url == "" { return }
	b, _ := json.Marshal(dec)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil { log.Printf("post cti: %v", err); return }
	io.Copy(io.Discard, resp.Body); resp.Body.Close()
}

func (s *Server) HandleStream(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil { log.Printf("ws upgrade: %v", err); return }
	defer c.Close()
	mActiveSessions.Inc()
	defer mActiveSessions.Dec()

	mode := os.Getenv("ASR_MODE") // shadow | force
	if mode == "" { mode = "shadow" }

	agg := strings.Builder{}
	tryDecide := func(text string) {
		phr := s.phrases.Load().(*asr.Phrases)
		if cls, ok := phr.Match(text); ok {
			el := time.Since(start).Milliseconds()
			mDecisions.WithLabelValues(cls).Inc()
			mLatency.Observe(float64(el))
			dec := Decision{UUID: r.URL.Query().Get("uuid"), Result: cls, Confidence: 0.9, LatencyMs: int(el), Transcript: text, Mode: "early"}
			if mode == "force" {
				s.callbackCTI(dec)
			}
		}
	}

	fasrURL := os.Getenv("FUNASR_WS_URL")
	var fasr *fa.Client
	var transcriptHandler = func(text string) {
		agg.WriteString(text)
		// strong phrase match triggers immediately
		tryDecide(text)
		// or aggregated content over ~1s
		if time.Since(start) > 1200*time.Millisecond {
			tryDecide(agg.String())
		}
	}

	if fasrURL != "" {
		fasr = fa.New(fasrURL)
		fasr.OnTranscript = transcriptHandler
		if err := fasr.Connect(); err != nil { log.Printf("funasr connect: %v", err) } else { go fasr.ReadLoop() }
		defer func(){ if fasr != nil { fasr.Close() } }()
	}

	// uuid used via r.URL.Query().Get("uuid") in decision path
	for {
		mt, data, err := c.ReadMessage()
		if err != nil { if err == io.EOF { return }; log.Printf("ws read: %v", err); return }
		switch mt {
		case websocket.TextMessage:
			text := string(data)
			agg.WriteString(text)
			tryDecide(text)
			if time.Since(start) > 1200*time.Millisecond {
				tryDecide(agg.String())
			}
		case websocket.BinaryMessage:
			if fasr != nil {
				_ = fasr.SendPCM(data)
				mForwardedFrames.Inc()
			}
		}
	}
}