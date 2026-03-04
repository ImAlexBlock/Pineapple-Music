package service

import (
	"testing"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=ON"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	db.AutoMigrate(
		&model.AdminKey{}, &model.GuestKey{}, &model.Session{},
		&model.Setting{}, &model.Track{}, &model.TrackLyric{},
		&model.Playlist{}, &model.PlaylistTrack{},
		&model.ScanJob{}, &model.ScanConflict{},
		&model.PlayEvent{}, &model.AuditLog{},
	)
	return db
}

func testConfig() *config.Config {
	return &config.Config{
		Port:          3880,
		DataDir:       "./test-data",
		SessionMaxAge: 3600,
	}
}

func TestBootstrapAndLogin(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	svc := &AuthService{DB: db, Cfg: cfg}

	// Not bootstrapped initially
	if svc.IsBootstrapped() {
		t.Fatal("expected not bootstrapped")
	}

	// Bootstrap
	key, err := svc.Bootstrap()
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if key == "" {
		t.Fatal("expected non-empty key")
	}

	// Now bootstrapped
	if !svc.IsBootstrapped() {
		t.Fatal("expected bootstrapped")
	}

	// Cannot bootstrap again
	_, err = svc.Bootstrap()
	if err == nil {
		t.Fatal("expected error on double bootstrap")
	}

	// Login with admin key
	session, role, err := svc.Login(key, "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if role != "admin" {
		t.Fatalf("expected admin role, got %s", role)
	}
	if session.ID == "" {
		t.Fatal("expected session ID")
	}

	// Validate session
	sess, err := svc.ValidateSession(session.ID)
	if err != nil {
		t.Fatalf("validate session: %v", err)
	}
	if sess.Role != "admin" {
		t.Fatalf("expected admin role in session")
	}

	// Guest key: plaintext is no longer stored in DB.
	// Test guest login by rotating the guest key (which returns the new plaintext).
	guestKey, err := svc.RotateGuestKey()
	if err != nil {
		t.Fatalf("rotate guest key: %v", err)
	}
	_, guestRole, err := svc.Login(guestKey, "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("guest login: %v", err)
	}
	if guestRole != "guest" {
		t.Fatalf("expected guest role, got %s", guestRole)
	}

	// Invalid key
	_, _, err = svc.Login("invalid-key", "127.0.0.1", "test")
	if err == nil {
		t.Fatal("expected error for invalid key")
	}

	// Logout
	svc.Logout(session.ID)
	_, err = svc.ValidateSession(session.ID)
	if err == nil {
		t.Fatal("expected error after logout")
	}
}

func TestRotateKeys(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	svc := &AuthService{DB: db, Cfg: cfg}

	key, _ := svc.Bootstrap()

	// Rotate admin key
	newKey, err := svc.RotateAdminKey()
	if err != nil {
		t.Fatalf("rotate admin key: %v", err)
	}
	if newKey == key {
		t.Fatal("expected different key")
	}

	// Old key should not work
	if svc.ValidateAdminKey(key) {
		t.Fatal("old key should be invalid")
	}

	// New key should work
	if !svc.ValidateAdminKey(newKey) {
		t.Fatal("new key should be valid")
	}
}
