package service

import (
	"testing"
)

func TestPlaylistCRUD(t *testing.T) {
	db := setupTestDB(t)
	svc := &PlaylistService{DB: db}

	// Create
	pl, err := svc.Create("Test Playlist")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if pl.Name != "Test Playlist" {
		t.Fatalf("expected name 'Test Playlist', got %s", pl.Name)
	}

	// List
	playlists, err := svc.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(playlists) != 1 {
		t.Fatalf("expected 1 playlist, got %d", len(playlists))
	}

	// Get
	got, err := svc.Get(pl.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "Test Playlist" {
		t.Fatalf("expected name 'Test Playlist', got %s", got.Name)
	}

	// Delete
	err = svc.Delete(pl.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	playlists, _ = svc.List()
	if len(playlists) != 0 {
		t.Fatalf("expected 0 playlists, got %d", len(playlists))
	}
}
