package asr_gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/cati/system/internal/asr"
)

var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

type Server struct {
	phrases atomic.Value // *asr.Phrases
}

type Decision struct {
	UUID        string  `json:"uuid"`
	Result      string  `json:"result"`
	Confidence  float64 `json:"confidence"`
	LatencyMs   int     `json:"latency_ms"`
	Transcript  string  `json:"transcript"`
	Mode        string  `json:"mode"`
	ProofURI    string  `json:"audio_proof_uri"`
	Fallback    bool    `json:"fallback"`
}

func NewServer() (*Server, error) {
	p := &Server{}
	if err := p.loadPhrases(); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Server) loadPhrases() error {
	path := os.Getenv("PHRASES_PATH")
	if path == "" {
		path = filepath.Join("/", "config", "phrases.yml")
	}
	phr, err := asr.LoadPhrases(path)
	if err != nil {
		return fmt.Errorf("load phrases: %w", err)
	}
	s.phrases.Store(phr)
	return nil
}

func (s *Server) HandleReload(w http.ResponseWriter, r *http.Request) {
	if err := s.loadPhrases(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reloaded"))
}

// HandleStream accepts WS but does not process PCM here; placeholder for FunASR integration.
func (s *Server) HandleStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}
	defer conn.Close()

	uuid := r.URL.Query().Get("uuid")
	_ = uuid

	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("ws read: %v", err)
			return
		}
		// Placeholder: accept text message as transcript for initial testing
		if mt == websocket.TextMessage {
			text := string(data)
			phr := s.phrases.Load().(*asr.Phrases)
			if cls, ok := phr.Match(text); ok {
				dec := Decision{UUID: uuid, Result: cls, Confidence: 0.99, LatencyMs: 200, Transcript: text, Mode: "early"}
				b, _ := json.Marshal(dec)
				// Echo back as decision frame (in production, POST to CTI)
				_ = conn.WriteMessage(websocket.TextMessage, b)
			}
		}
	}
}