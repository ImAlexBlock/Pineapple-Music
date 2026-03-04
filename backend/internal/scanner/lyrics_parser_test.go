package scanner

import (
	"testing"
)

func TestParseLRC(t *testing.T) {
	lrc := `[ti:Test Song]
[ar:Test Artist]
[00:01.00]First line
[00:05.50]Second line
[00:10.00]Third line`

	lines, synced := ParseLRC(lrc)
	if !synced {
		t.Fatal("expected synced")
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	// Check first line
	if lines[0].Time != 1.0 {
		t.Fatalf("expected time 1.0, got %f", lines[0].Time)
	}
	if lines[0].Text != "First line" {
		t.Fatalf("expected 'First line', got '%s'", lines[0].Text)
	}

	// Check second line
	if lines[1].Time != 5.5 {
		t.Fatalf("expected time 5.5, got %f", lines[1].Time)
	}
}

func TestParseLRC_NotSynced(t *testing.T) {
	plain := "Just some lyrics\nNo time tags here"
	_, synced := ParseLRC(plain)
	if synced {
		t.Fatal("expected not synced")
	}
}
