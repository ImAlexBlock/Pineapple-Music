package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"
	"pineapple-music/internal/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestRouter(t *testing.T) (*gorm.DB, http.Handler) {
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

	cfg := &config.Config{
		Port:           3880,
		DataDir:        "./test-data",
		SessionMaxAge:  3600,
		RateLimitRPS:   100,
		RateLimitBurst: 200,
		MaxUploadSize:  52428800,
	}

	r := SetupRouter(db, cfg)
	return db, r
}

func TestHealthCheck(t *testing.T) {
	_, handler := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}

func TestBootstrapFlow(t *testing.T) {
	db, handler := setupTestRouter(t)
	cfg := &config.Config{
		Port:           3880,
		DataDir:        "./test-data",
		SessionMaxAge:  3600,
		RateLimitRPS:   100,
		RateLimitBurst: 200,
		MaxUploadSize:  52428800,
	}

	// Check status - not bootstrapped
	req := httptest.NewRequest("GET", "/api/v1/setup/status", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var status map[string]bool
	json.Unmarshal(w.Body.Bytes(), &status)
	if status["bootstrapped"] {
		t.Fatal("expected not bootstrapped")
	}

	// Bootstrap via service directly (HTTP bootstrap no longer returns the key)
	authSvc := &service.AuthService{DB: db, Cfg: cfg}
	adminKey, err := authSvc.Bootstrap()
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	// Get CSRF token
	csrfCookie := ""
	for _, c := range w.Result().Cookies() {
		if c.Name == "pm_csrf" {
			csrfCookie = c.Value
		}
	}

	// Login
	loginBody := `{"key":"` + adminKey + `"}`
	req = httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	if csrfCookie != "" {
		req.AddCookie(&http.Cookie{Name: "pm_csrf", Value: csrfCookie})
		req.Header.Set("X-CSRF-Token", csrfCookie)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("login expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Extract session cookie
	var sessionCookie string
	for _, c := range w.Result().Cookies() {
		if c.Name == "pm_session" {
			sessionCookie = c.Value
		}
	}
	if sessionCookie == "" {
		t.Fatal("expected session cookie")
	}

	// Access /auth/me
	req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "pm_session", Value: sessionCookie})
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("me expected 200, got %d", w.Code)
	}

	var meResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &meResp)
	if meResp["role"] != "admin" {
		t.Fatalf("expected admin role, got %s", meResp["role"])
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	_, handler := setupTestRouter(t)

	// Access protected route without session
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
