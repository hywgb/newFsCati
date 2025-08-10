package cticontroller

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Store struct { db *sql.DB }

func NewStoreFromEnv() (*Store, error) {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" { dsn = "postgres://user:pass@postgres:5432/cati?sslmode=disable" }
	db, err := sql.Open("postgres", dsn)
	if err != nil { return nil, err }
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &Store{db: db}, nil
}

func (s *Store) SaveAsrDecision(d AsrDecision) error {
	_, err := s.db.Exec(`
		insert into call_attempts (uuid, asr_result, asr_confidence, asr_latency_ms, asr_transcript, asr_mode)
		values ($1,$2,$3,$4,$5,$6)
		on conflict (uuid) do update set asr_result=excluded.asr_result, asr_confidence=excluded.asr_confidence, asr_latency_ms=excluded.asr_latency_ms, asr_transcript=excluded.asr_transcript, asr_mode=excluded.asr_mode
	`, d.UUID, d.Result, d.Confidence, d.LatencyMs, d.Transcript, d.Mode)
	return err
}

func (s *Store) Close() { if s.db != nil { _ = s.db.Close() } }