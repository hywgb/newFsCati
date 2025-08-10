package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cati/system/services/asr-gateway"
)

func main() {
	addr := ":10000"
	if v := os.Getenv("ASR_GATEWAY_ADDR"); v != "" {
		addr = v
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server, err := asr_gateway.NewServer()
	if err != nil {
		log.Fatalf("init asr gateway: %v", err)
	}
	http.HandleFunc("/stream", server.HandleStream)
	http.HandleFunc("/config/reload", server.HandleReload)

	log.Printf("ASR Gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("serve: %v", err)
	}
}