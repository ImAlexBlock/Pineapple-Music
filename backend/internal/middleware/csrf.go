package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

const (
	CSRFCookie = "pm_csrf"
	CSRFHeader = "X-CSRF-Token"
)

// CSRF implements double-submit cookie pattern.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Safe methods don't need CSRF check
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			// Ensure CSRF cookie exists
			if _, err := c.Cookie(CSRFCookie); err != nil {
				token := generateCSRFToken()
				c.SetCookie(CSRFCookie, token, 86400, "/", "", false, false) // not HttpOnly so JS can read
			}
			c.Next()
			return
		}

		// For state-changing methods, validate CSRF
		cookieToken, err := c.Cookie(CSRFCookie)
		if err != nil || cookieToken == "" {
			util.ErrorResponse(c, http.StatusForbidden, "csrf_missing", "CSRF cookie missing")
			c.Abort()
			return
		}

		headerToken := c.GetHeader(CSRFHeader)
		if headerToken == "" || headerToken != cookieToken {
			util.ErrorResponse(c, http.StatusForbidden, "csrf_invalid", "CSRF token mismatch")
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateCSRFToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
