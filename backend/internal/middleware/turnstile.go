package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

const turnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type turnstileResponse struct {
	Success bool `json:"success"`
}

// Turnstile validates Cloudflare Turnstile token if secret is configured.
func Turnstile(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.Next()
			return
		}

		token := c.GetHeader("X-Turnstile-Token")
		if token == "" {
			token = c.PostForm("turnstile_token")
		}
		if token == "" {
			util.ErrorResponse(c, http.StatusBadRequest, "turnstile_missing", "Turnstile token required")
			c.Abort()
			return
		}

		if !verifyTurnstile(secret, token, c.ClientIP()) {
			util.ErrorResponse(c, http.StatusForbidden, "turnstile_failed", "Turnstile verification failed")
			c.Abort()
			return
		}

		c.Next()
	}
}

func verifyTurnstile(secret, token, ip string) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(turnstileVerifyURL, map[string][]string{
		"secret":   {secret},
		"response": {token},
		"remoteip": {ip},
	})
	if err != nil {
		fmt.Printf("Turnstile verify error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	var result turnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}
	return result.Success
}
