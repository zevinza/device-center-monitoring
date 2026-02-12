package middleware

import (
	"api/utils/resp"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config holds middleware configuration
type Config struct {
	JWTSecret           string
	APISecretKey        string
	TokenHeader         string
	TokenPrefix         string
	SecretKeyHeader     string
	UseRedisCache       bool
	CacheTTL            time.Duration
	CheckPermissions    bool
	RequiredRoles       []string
	RequiredPermissions []string
	CheckMethod         string // "all", "write", "read"
}

// Option is a function that configures the middleware
type Option func(*Config)

// WithJWTSecret sets the JWT secret key
func WithJWTSecret(secret string) Option {
	return func(c *Config) {
		c.JWTSecret = secret
	}
}

// WithAPISecretKey sets the API secret key for public API authentication
func WithAPISecretKey(secretKey string) Option {
	return func(c *Config) {
		c.APISecretKey = secretKey
	}
}

// WithSecretKeyHeader sets the header name for API secret key
func WithSecretKeyHeader(header string) Option {
	return func(c *Config) {
		c.SecretKeyHeader = header
	}
}

// WithTokenHeader sets the token header name
func WithTokenHeader(header string) Option {
	return func(c *Config) {
		c.TokenHeader = header
	}
}

// WithRedisCache enables/disables Redis caching
func WithRedisCache(enable bool) Option {
	return func(c *Config) {
		c.UseRedisCache = enable
	}
}

// WithCacheTTL sets the cache TTL
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.CacheTTL = ttl
	}
}

// WithRequiredRoles sets required roles for authorization
func WithRequiredRoles(roles ...string) Option {
	return func(c *Config) {
		c.CheckPermissions = true
		c.RequiredRoles = roles
	}
}

// WithRequiredPermissions sets required permissions for authorization
func WithRequiredPermissions(permissions ...string) Option {
	return func(c *Config) {
		c.CheckPermissions = true
		c.RequiredPermissions = permissions
	}
}

// WithCheckMethod sets when to check permissions (all/write/read)
func WithCheckMethod(method string) Option {
	return func(c *Config) {
		c.CheckMethod = method
	}
}

// RequireRole is a convenience function that returns middleware requiring specific roles
func RequireRole(roles ...string) fiber.Handler {
	return OAuth2Authentication(WithRequiredRoles(roles...))
}

// RequirePermission is a convenience function that returns middleware requiring specific permissions
func RequirePermission(permissions ...string) fiber.Handler {
	return OAuth2Authentication(WithRequiredPermissions(permissions...))
}

// RequireAny is a convenience function that requires either role OR permission
func RequireAny(roles []string, permissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First authenticate
		if err := OAuth2Authentication()(c); err != nil {
			return err
		}

		// Check if user has any required role
		userRoles, _ := GetRoles(c)
		for _, userRole := range userRoles {
			for _, requiredRole := range roles {
				if userRole == requiredRole {
					return c.Next()
				}
			}
		}

		// Check if user has any required permission
		userPermissions, _ := GetPermissions(c)
		for _, userPerm := range userPermissions {
			for _, requiredPerm := range permissions {
				if userPerm == requiredPerm {
					return c.Next()
				}
			}
		}

		return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient permissions"))
	}
}

// // ErrorForbidden returns a forbidden error response
// func ErrorForbidden(c *fiber.Ctx, message string) error {
// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 		"code":    fiber.StatusForbidden,
// 		"status":  false,
// 		"message": message,
// 	})
// }
