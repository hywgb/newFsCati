package cticontroller

import (
	"encoding/json"
	"net/http"
)

type AsrDecision struct {
	UUID        string  `json:"uuid"`
	Result      string  `json:"result"`
	Confidence  float64 `json:"confidence"`
	LatencyMs   int     `json:"latency_ms"`
	Transcript  string  `json:"transcript"`
	Mode        string  `json:"mode"`
	ProofURI    string  `json:"audio_proof_uri"`
	Fallback    bool    `json:"fallback"`
}

func (s *Server) HandleAsrDecision(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var d AsrDecision
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mAsrCallbacks.WithLabelValues(d.Result).Inc()
	if s.store != nil {
		_ = s.store.SaveAsrDecision(d)
	}
	if d.Confidence >= 0.75 {
		s.KillByDecision(d.UUID)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}