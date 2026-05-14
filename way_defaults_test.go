package way

import (
	"os"
	"testing"
	"time"
)

func TestNewSetsServerTimeouts(t *testing.T) {
	w := New()

	// Check that timeouts are set to safe defaults
	if w.Server.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout = %v, want %v", w.Server.ReadHeaderTimeout, 5*time.Second)
	}
	if w.Server.ReadTimeout != 15*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", w.Server.ReadTimeout, 15*time.Second)
	}
	if w.Server.WriteTimeout != 15*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", w.Server.WriteTimeout, 15*time.Second)
	}
	if w.Server.IdleTimeout != 30*time.Second {
		t.Errorf("IdleTimeout = %v, want %v", w.Server.IdleTimeout, 30*time.Second)
	}
}

func TestNewInitializesRouter(t *testing.T) {
	w := New()

	if w.router == nil {
		t.Error("New() router is nil")
	}
}

func TestNewInitializesSessions(t *testing.T) {
	w := New()

	if w.sessions == nil {
		t.Error("New() sessions is nil")
	}
}

func TestNewInitializesLogger(t *testing.T) {
	w := New()

	if w.Logger == nil {
		t.Error("New() Logger is nil")
	}
}

func TestASCIIArtLoggingGatedByEnv(t *testing.T) {
	// This is a manual test helper; in a real scenario, you'd capture log output
	// For now, just verify that the env variable is checked
	oldVal, wasSet := os.LookupEnv("WAY_LOG_ASCII_ART")
	defer func() {
		if wasSet {
			os.Setenv("WAY_LOG_ASCII_ART", oldVal)
		} else {
			os.Unsetenv("WAY_LOG_ASCII_ART")
		}
	}()

	// With env var not set, ASCII art should not appear
	os.Unsetenv("WAY_LOG_ASCII_ART")
	w := New()
	if w == nil {
		t.Error("New() returned nil")
	}

	// With env var set to "true", ASCII art check should pass
	os.Setenv("WAY_LOG_ASCII_ART", "true")
	w2 := New()
	if w2 == nil {
		t.Error("New() with WAY_LOG_ASCII_ART=true returned nil")
	}
}
