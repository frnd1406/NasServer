package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/services"
)

// VaultGuard blocks requests when the vault is locked
// This middleware should be applied to routes that require encryption access
func VaultGuard(encSvc *services.EncryptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !encSvc.IsUnlocked() {
			c.JSON(http.StatusLocked, gin.H{
				"error":   "vault is locked",
				"message": "Please unlock the vault to access this resource",
				"code":    "VAULT_LOCKED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// VaultConfigured ensures the vault is configured before proceeding
// Used for routes that require the vault to be set up first
func VaultConfigured(encSvc *services.EncryptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !encSvc.IsConfigured() {
			c.JSON(http.StatusPreconditionFailed, gin.H{
				"error":   "vault not configured",
				"message": "Please run vault setup first",
				"code":    "VAULT_NOT_CONFIGURED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
