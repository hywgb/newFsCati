package main

import (
	"log"
	"net/http"
	"os"

	cti "github.com/cati/system/services/cti-controller"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("CTI_HTTP_ADDR"); v != "" {
		addr = v
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/asr/decision", cti.HandleAsrDecision)

	log.Printf("CTI Controller HTTP listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("serve: %v", err)
	}
}