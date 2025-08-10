package cticontroller

import (
	"encoding/json"
	"log"
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

func HandleAsrDecision(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var d AsrDecision
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: write to DB and send uuid_kill via ESL client (placeholder log)
	log.Printf("ASR Decision: uuid=%s result=%s conf=%.2f latency=%dms", d.UUID, d.Result, d.Confidence, d.LatencyMs)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}