package asr

import (
	"path/filepath"
	"runtime"
	"testing"
)

func phrasesPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "config", "phrases.yml")
}

func TestMatchPowerOff(t *testing.T) {
	p, err := LoadPhrases(phrasesPath())
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cls, ok := p.Match("您拨打的用户已关机"); !ok || cls != "power_off" {
		t.Fatalf("expected power_off got %v ok=%v", cls, ok)
	}
}

func TestMatchInvalidNumber(t *testing.T) {
	p, err := LoadPhrases(phrasesPath())
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cls, ok := p.Match("您拨打的号码是空号"); !ok || cls != "invalid_number" {
		t.Fatalf("expected invalid_number got %v ok=%v", cls, ok)
	}
}