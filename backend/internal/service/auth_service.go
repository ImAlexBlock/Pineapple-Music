package service

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"gorm.io/gorm"
)

type AuthService struct {
	DB  *gorm.DB
	Cfg *config.Config
}

// IsBootstrapped checks if an admin key exists.
func (s *AuthService) IsBootstrapped() bool {
	var count int64
	s.DB.Model(&model.AdminKey{}).Count(&count)
	return count > 0
}

// Bootstrap creates the initial admin key. Returns the plaintext key.
func (s *AuthService) Bootstrap() (string, error) {
	if s.IsBootstrapped() {
		return "", fmt.Errorf("already bootstrapped")
	}

	key, err := util.GenerateRandomKey(32)
	if err != nil {
		return "", err
	}

	adminKey := model.AdminKey{
		KeyHash: util.HashKey(key),
		MD5Hash: md5Hex(key),
	}
	if err := s.DB.Create(&adminKey).Error; err != nil {
		return "", err
	}

	// Also create a default guest key — only hashes stored
	guestPlain, err := util.GenerateRandomKey(16)
	if err != nil {
		return "", err
	}
	guestKey := model.GuestKey{
		KeyHash: util.HashKey(guestPlain),
		MD5Hash: md5Hex(guestPlain),
	}
	s.DB.Create(&guestKey)

	// Print guest key to server log (NOT stored as plaintext in DB)
	fmt.Printf("  Guest key: %s\n", guestPlain)

	return key, nil
}

// ValidateAdminKey checks if the given key matches any admin key.
func (s *AuthService) ValidateAdminKey(key string) bool {
	hash := util.HashKey(key)
	var count int64
	s.DB.Model(&model.AdminKey{}).Where("key_hash = ?", hash).Count(&count)
	return count > 0
}

// ValidateGuestKey checks if the given key matches any guest key.
func (s *AuthService) ValidateGuestKey(key string) bool {
	hash := util.HashKey(key)
	var count int64
	s.DB.Model(&model.GuestKey{}).Where("key_hash = ?", hash).Count(&count)
	return count > 0
}

// ValidateSubsonicToken validates Subsonic token auth: token = md5(password + salt).
// We store md5(password) as MD5Hash. Subsonic token = md5(password + salt).
// Since md5 is not composable, we cannot derive token from md5(password) alone.
// However, we store the MD5 of the *key* itself (md5Hex(key)).
// Subsonic clients typically send token = md5(password + salt) where password is the key.
// We iterate over stored keys, compute md5(key + salt) using the stored MD5 hash approach.
// NOTE: This requires us to store enough info to reconstruct the token.
// Since we can't reverse the hash, we use a pragmatic approach: store md5(key) and
// try token = md5(md5(key) + salt) — but standard Subsonic uses md5(plaintext + salt).
// This means true token auth is impossible without plaintext. We document this limitation
// and fall through to password auth.
// UPDATE: We now store md5(plaintext_key) in MD5Hash. Many Subsonic clients also support
// sending the password as md5(password) already. We check both approaches:
// 1. Standard: token = md5(plaintext + salt) — cannot verify without plaintext
// 2. Alternative: some clients send token = md5(hex_md5_password + salt) — we can verify this
func (s *AuthService) ValidateSubsonicToken(token, salt string) (string, bool) {
	// Check admin keys
	var adminKeys []model.AdminKey
	s.DB.Where("md5_hash != ''").Find(&adminKeys)
	for _, ak := range adminKeys {
		// Try: token = md5(md5hex(password) + salt) — some clients use this
		expected := md5Hex(ak.MD5Hash + salt)
		if expected == token {
			return "admin", true
		}
	}

	// Check guest keys
	var guestKeys []model.GuestKey
	s.DB.Where("md5_hash != ''").Find(&guestKeys)
	for _, gk := range guestKeys {
		expected := md5Hex(gk.MD5Hash + salt)
		if expected == token {
			return "guest", true
		}
	}

	return "", false
}

// Login validates credentials and creates a session.
func (s *AuthService) Login(key, ip, userAgent string) (*model.Session, string, error) {
	var role string
	if s.ValidateAdminKey(key) {
		role = "admin"
	} else if s.ValidateGuestKey(key) {
		role = "guest"
	} else {
		return nil, "", fmt.Errorf("invalid key")
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, "", err
	}

	session := &model.Session{
		ID:        sessionID,
		Role:      role,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(time.Duration(s.Cfg.SessionMaxAge) * time.Second),
	}

	if err := s.DB.Create(session).Error; err != nil {
		return nil, "", err
	}

	return session, role, nil
}

// ValidateSession checks if a session is valid and optionally extends it.
func (s *AuthService) ValidateSession(sessionID string) (*model.Session, error) {
	var session model.Session
	if err := s.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, fmt.Errorf("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		s.DB.Delete(&session)
		return nil, fmt.Errorf("session expired")
	}

	// Sliding expiration: extend if more than half the time has passed
	halfLife := time.Duration(s.Cfg.SessionMaxAge) * time.Second / 2
	if time.Until(session.ExpiresAt) < halfLife {
		session.ExpiresAt = time.Now().Add(time.Duration(s.Cfg.SessionMaxAge) * time.Second)
		s.DB.Save(&session)
	}

	return &session, nil
}

// Logout deletes a session.
func (s *AuthService) Logout(sessionID string) {
	s.DB.Delete(&model.Session{}, "id = ?", sessionID)
}

// CleanExpiredSessions removes all expired sessions.
func (s *AuthService) CleanExpiredSessions() {
	s.DB.Where("expires_at < ?", time.Now()).Delete(&model.Session{})
}

// RotateAdminKey replaces all admin keys. Returns the new plaintext key (show once).
func (s *AuthService) RotateAdminKey() (string, error) {
	key, err := util.GenerateRandomKey(32)
	if err != nil {
		return "", err
	}

	tx := s.DB.Begin()
	tx.Where("1=1").Delete(&model.AdminKey{})
	tx.Create(&model.AdminKey{
		KeyHash: util.HashKey(key),
		MD5Hash: md5Hex(key),
	})
	tx.Where("role = ?", "admin").Delete(&model.Session{})
	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return key, nil
}

// RotateGuestKey replaces all guest keys. Returns the new plaintext key (show once).
func (s *AuthService) RotateGuestKey() (string, error) {
	key, err := util.GenerateRandomKey(16)
	if err != nil {
		return "", err
	}

	tx := s.DB.Begin()
	tx.Where("1=1").Delete(&model.GuestKey{})
	tx.Create(&model.GuestKey{
		KeyHash: util.HashKey(key),
		MD5Hash: md5Hex(key),
	})
	tx.Where("role = ?", "guest").Delete(&model.Session{})
	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return key, nil
}

// GetSetting retrieves a setting value.
func (s *AuthService) GetSetting(key string) string {
	var setting model.Setting
	if err := s.DB.First(&setting, "key = ?", key).Error; err != nil {
		return ""
	}
	return setting.Value
}

func md5Hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
