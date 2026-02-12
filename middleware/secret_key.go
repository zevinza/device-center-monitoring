package middleware

import (
	"api/utils/resp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

const (
	AuthModeKey    = "auth_mode"
	AuthModeSecret = "secret_key"
	AuthModeJWT    = "jwt"
	SecretKeyKey   = "api_secret_key"
)

// SecretKeyAuthentication validates API secret key for public/non-auth APIs
func SecretKeyAuthentication(options ...Option) fiber.Handler {
	config := &Config{
		APISecretKey:    viper.GetString("API_SECRET_KEY"),
		SecretKeyHeader: "X-API-Key", // Default header name
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	return func(c *fiber.Ctx) error {
		// Extract secret key from header
		secretKey := c.Get(config.SecretKeyHeader)
		if secretKey == "" {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("API secret key required"))
		}

		// Validate secret key
		if secretKey != config.APISecretKey {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Invalid API secret key"))
		}

		// Store authentication mode in context
		c.Locals(AuthModeKey, AuthModeSecret)
		c.Locals(SecretKeyKey, secretKey)

		return c.Next()
	}
}

// DualAuthentication supports both secret key (for public APIs) and JWT (for authenticated users)
// It tries secret key first, then falls back to JWT if secret key is not present
func DualAuthentication(options ...Option) fiber.Handler {
	config := &Config{
		JWTSecret:           viper.GetString("JWT_SECRET"),
		APISecretKey:        viper.GetString("API_SECRET_KEY"),
		TokenHeader:         "Authorization",
		TokenPrefix:         "Bearer ",
		SecretKeyHeader:     "X-API-Key",
		UseRedisCache:       viper.GetBool("USE_REDIS_CACHE"),
		CacheTTL:            10 * time.Minute,
		CheckPermissions:    false,
		RequiredRoles:       []string{},
		RequiredPermissions: []string{},
		CheckMethod:         "all",
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	return func(c *fiber.Ctx) error {
		// Try secret key authentication first (for public APIs)
		secretKey := c.Get(config.SecretKeyHeader)
		if secretKey != "" {
			if secretKey == config.APISecretKey {
				// Valid secret key - public API access
				c.Locals(AuthModeKey, AuthModeSecret)
				c.Locals(SecretKeyKey, secretKey)
				return c.Next()
			}
			// Invalid secret key
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Invalid API secret key"))
		}

		// No secret key provided, try JWT authentication
		authHeader := c.Get(config.TokenHeader)
		if authHeader == "" {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Authorization required: provide either API secret key or JWT token"))
		}

		if !strings.HasPrefix(authHeader, config.TokenPrefix) {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Invalid authorization header format"))
		}

		tokenString := strings.TrimPrefix(authHeader, config.TokenPrefix)
		if tokenString == "" {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Token not found"))
		}

		// Validate JWT token
		claims, err := validateToken(tokenString, config.JWTSecret)
		if err != nil {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Invalid or expired token"))
		}

		// Store JWT authentication mode
		c.Locals(AuthModeKey, AuthModeJWT)

		// Handle permissions if needed
		if config.CheckPermissions && shouldCheckPermissions(c, config) {
			permissions, err := getFreshPermissions(c.Context(), claims.UserID, claims, config)
			if err != nil {
				storeClaims(c, claims.UserID, claims.Email, claims.Roles, claims.Permissions)
			} else {
				storeClaims(c, claims.UserID, claims.Email, permissions.Roles, permissions.Permissions)
			}
		} else {
			storeClaims(c, claims.UserID, claims.Email, claims.Roles, claims.Permissions)
		}

		// Check authorization if required
		if config.CheckPermissions {
			if len(config.RequiredRoles) > 0 {
				if !hasRole(c, config.RequiredRoles) {
					return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient role permissions"))
				}
			}

			if len(config.RequiredPermissions) > 0 {
				if !hasPermission(c, config.RequiredPermissions) {
					return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient permissions"))
				}
			}
		}

		return c.Next()
	}
}

// GetAuthMode returns the authentication mode used (secret_key or jwt)
func GetAuthMode(c *fiber.Ctx) (string, bool) {
	mode, ok := c.Locals(AuthModeKey).(string)
	return mode, ok
}

// IsSecretKeyAuth checks if the request was authenticated using secret key
func IsSecretKeyAuth(c *fiber.Ctx) bool {
	mode, ok := c.Locals(AuthModeKey).(string)
	return ok && mode == AuthModeSecret
}

// IsJWTAuth checks if the request was authenticated using JWT
func IsJWTAuth(c *fiber.Ctx) bool {
	mode, ok := c.Locals(AuthModeKey).(string)
	return ok && mode == AuthModeJWT
}
