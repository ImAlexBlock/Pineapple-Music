package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"pineapple-music/internal/config"
	"pineapple-music/internal/middleware"
	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

func BootstrapStatus(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"bootstrapped": authSvc.IsBootstrapped(),
		})
	}
}

func Bootstrap(authSvc *service.AuthService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := authSvc.Bootstrap()
		if err != nil {
			util.ErrorResponse(c, http.StatusConflict, "already_bootstrapped", err.Error())
			return
		}

		auditSvc.Log("bootstrap", "admin", c.ClientIP(), "System bootstrapped")

		// Print admin key to server console only
		log.Printf("=== ADMIN KEY (save this, shown only once): %s ===", key)
		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Printf("║  ADMIN KEY: %-40s ║\n", key)
		fmt.Println("║  Save this key! It will NOT be shown again.         ║")
		fmt.Println("╚══════════════════════════════════════════════════════╝")
		fmt.Println()

		// Do NOT return the key over HTTP
		c.JSON(http.StatusOK, gin.H{
			"message": "System initialized. The admin key has been printed to the server console.",
		})
	}
}

func Login(authSvc *service.AuthService, cfg *config.Config, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Key string `json:"key" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "Key is required")
			return
		}

		session, role, err := authSvc.Login(req.Key, c.ClientIP(), c.GetHeader("User-Agent"))
		if err != nil {
			auditSvc.Log("login_failed", "", c.ClientIP(), "Invalid key attempt")
			util.ErrorResponse(c, http.StatusUnauthorized, "invalid_key", "Invalid key")
			return
		}

		// Determine Secure flag: explicit config or auto-detect from X-Forwarded-Proto
		secure := cfg.SecureCookie
		if !secure && strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
			secure = true
		}

		// SameSite must be set BEFORE SetCookie
		c.SetSameSite(http.SameSiteLaxMode)

		c.SetCookie(
			middleware.SessionCookie,
			session.ID,
			int(time.Until(session.ExpiresAt).Seconds()),
			"/",
			"",
			secure,
			true, // httpOnly
		)

		auditSvc.Log("login", role, c.ClientIP(), "")

		c.JSON(http.StatusOK, gin.H{
			"role": role,
		})
	}
}

func Logout(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, _ := c.Cookie(middleware.SessionCookie)
		if sessionID != "" {
			authSvc.Logout(sessionID)
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(middleware.SessionCookie, "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}

func Me(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(middleware.ContextRole)
		c.JSON(http.StatusOK, gin.H{
			"role": role,
		})
	}
}
