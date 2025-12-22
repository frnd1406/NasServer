package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/repository"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware validates JWT tokens and checks blacklist
func AuthMiddleware(jwtService *services.JWTService, redis *database.RedisClient, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.WithField("request_id", requestID).Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "Missing authorization token",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.WithField("request_id", requestID).Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "Invalid authorization header format",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate JWT token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("Invalid JWT token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "Invalid or expired token",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Check if token is blacklisted (logout revocation)
		ctx := context.Background()
		blacklisted, err := redis.Get(ctx, "blacklist:"+tokenString).Result()
		if err == nil && blacklisted == "1" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"user_id":    claims.UserID,
			}).Warn("Blacklisted token used")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "Token has been revoked",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Store user info in context for handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    claims.UserID,
		}).Debug("User authenticated successfully")

		c.Next()
	}
}

// RequireRole creates middleware that enforces role-based access control
// IMPORTANT: Must be used AFTER AuthMiddleware
func RequireRole(userRepo *repository.UserRepository, requiredRole models.UserRole, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		userID := c.GetString("user_id")

		// Check if user_id exists in context (should be set by AuthMiddleware)
		if userID == "" {
			logger.WithField("request_id", requestID).Error("RBAC: user_id not found in context - AuthMiddleware not executed?")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "Authentication required",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Fetch user from database to check role
		ctx := c.Request.Context()
		user, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"user_id":    userID,
				"error":      err.Error(),
			}).Error("RBAC: failed to fetch user")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "internal_error",
					"message":    "Failed to verify permissions",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		if user == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"user_id":    userID,
			}).Warn("RBAC: user not found in database")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "User not found",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Check if user has required role
		if user.Role != requiredRole {
			logger.WithFields(logrus.Fields{
				"request_id":    requestID,
				"user_id":       userID,
				"user_role":     user.Role,
				"required_role": requiredRole,
			}).Warn("RBAC: access denied - insufficient permissions")
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":       "forbidden",
					"message":    "Insufficient permissions - admin access required",
					"request_id": requestID,
				},
			})
			c.Abort()
			return
		}

		// Store user role in context for downstream handlers
		c.Set("user_role", string(user.Role))

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
			"user_role":  user.Role,
		}).Debug("RBAC: access granted")

		c.Next()
	}
}

// AdminOnly is a convenience function for requiring admin role
func AdminOnly(userRepo *repository.UserRepository, logger *logrus.Logger) gin.HandlerFunc {
	return RequireRole(userRepo, models.RoleAdmin, logger)
}
