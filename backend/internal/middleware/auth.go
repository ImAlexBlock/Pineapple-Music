package middleware

import (
	"net"
	"net/http"
	"strings"

	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookie  = "pm_session"
	ContextRole    = "role"
	ContextSession = "session_id"
)

// Auth validates the session cookie and sets role in context.
func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookie)
		if err != nil || sessionID == "" {
			util.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", "Login required")
			c.Abort()
			return
		}

		session, err := authSvc.ValidateSession(sessionID)
		if err != nil {
			util.ErrorResponse(c, http.StatusUnauthorized, "session_invalid", err.Error())
			c.Abort()
			return
		}

		c.Set(ContextRole, session.Role)
		c.Set(ContextSession, session.ID)
		c.Next()
	}
}

// RequireAdmin ensures the user has admin role.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextRole)
		if role != "admin" {
			util.ErrorResponse(c, http.StatusForbidden, "forbidden", "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalAuth tries to authenticate but doesn't require it.
func OptionalAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookie)
		if err != nil || sessionID == "" {
			c.Next()
			return
		}

		session, err := authSvc.ValidateSession(sessionID)
		if err != nil {
			c.Next()
			return
		}

		c.Set(ContextRole, session.Role)
		c.Set(ContextSession, session.ID)
		c.Next()
	}
}

// AccessMode enforces private-mode access control.
// In private mode (access_mode != "public"), unauthenticated requests are rejected.
// Default (no setting stored) is "public" for ease of use.
func AccessMode(authSvc *service.AuthService, settingSvc *service.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mode := strings.TrimSpace(settingSvc.Get("access_mode"))
		// Default to public when no setting exists
		if mode == "" || strings.EqualFold(mode, "public") {
			c.Next()
			return
		}
		// "private" — require authentication (guest or admin)
		role, exists := c.Get(ContextRole)
		if !exists || role == nil || role == "" {
			util.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", "Login required (private mode)")
			c.Abort()
			return
		}
		c.Next()
	}
}

// LocalOnly restricts access to loopback addresses.
func LocalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !isLoopback(ip) {
			util.ErrorResponse(c, http.StatusForbidden, "forbidden", "This endpoint is only available from localhost")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isLoopback(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
