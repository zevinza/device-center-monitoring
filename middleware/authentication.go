package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"api/constant"
	"api/lib"
	"api/services/cache"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// CachedPermissions represents cached user permissions in Redis
type CachedPermissions struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Version     int64    `json:"version"` // Permission version for invalidation
	UpdatedAt   string   `json:"updated_at"`
}

const (
	ClaimsKey      = "claims"
	UserIDKey      = "user_id"
	EmailKey       = "email"
	RolesKey       = "roles"
	PermissionsKey = "permissions"
)

// OAuth2Authentication is a combined middleware that handles authentication and optional authorization
// It validates JWT tokens and optionally checks permissions/roles
func OAuth2Authentication(options ...Option) fiber.Handler {
	config := &Config{
		JWTSecret:           viper.GetString("JWT_SECRET"),
		TokenHeader:         "Authorization",
		TokenPrefix:         "Bearer ",
		UseRedisCache:       viper.GetBool("USE_REDIS_CACHE"),
		CacheTTL:            10 * time.Minute,
		CheckPermissions:    false,
		RequiredRoles:       []string{},
		RequiredPermissions: []string{},
		CheckMethod:         "all", // "all", "write", "read"
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	return func(c *fiber.Ctx) error {
		// 1. Extract and validate JWT token
		tokenString, err := extractToken(c, config)
		if err != nil {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized(err.Error()))
		}

		claims, err := validateToken(tokenString, config.JWTSecret)
		if err != nil {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Invalid or expired token"))
		}

		// 2. For GET requests: use token claims directly (fast, might be stale)
		// For write/admin operations: check Redis cache or DB for fresh permissions
		if config.CheckPermissions && shouldCheckPermissions(c, config) {
			permissions, err := getFreshPermissions(c.Context(), claims.UserID, claims, config)
			if err != nil {
				// Fallback to token claims if cache/DB fails
				storeClaims(c, claims.UserID, claims.Email, claims.Roles, claims.Permissions)
			} else {
				// Use fresh permissions from cache/DB
				storeClaims(c, claims.UserID, claims.Email, permissions.Roles, permissions.Permissions)
			}
		} else {
			// No permission check needed, just store token claims
			storeClaims(c, claims.UserID, claims.Email, claims.Roles, claims.Permissions)
		}

		// 3. Check authorization if required
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

// extractToken extracts JWT token from request header
func extractToken(c *fiber.Ctx, config *Config) (string, error) {
	authHeader := c.Get(config.TokenHeader)
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
	}

	if !strings.HasPrefix(authHeader, config.TokenPrefix) {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(authHeader, config.TokenPrefix)
	if tokenString == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Token not found")
	}

	return tokenString, nil
}

// validateToken validates JWT token and returns claims
func validateToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Token expired")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	return claims, nil
}

// getFreshPermissions gets fresh permissions from Redis cache or falls back to token claims
func getFreshPermissions(ctx context.Context, userID string, tokenClaims *TokenClaims, config *Config) (*CachedPermissions, error) {
	if !config.UseRedisCache || cache.Redis == nil {
		// Redis not available, return token claims
		return &CachedPermissions{
			Roles:       tokenClaims.Roles,
			Permissions: tokenClaims.Permissions,
		}, nil
	}

	cacheKey := "user:permissions:" + userID

	// Try to get from Redis
	cachedData, err := cache.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cached CachedPermissions
		if json.Unmarshal([]byte(cachedData), &cached) == nil {
			return &cached, nil
		}
	}

	// Cache miss: return token claims (in production, you'd fetch from DB here)
	// For now, we'll use token claims and optionally cache them
	permissions := &CachedPermissions{
		Roles:       tokenClaims.Roles,
		Permissions: tokenClaims.Permissions,
		Version:     time.Now().Unix(),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Cache the permissions for next time
	if config.UseRedisCache && cache.Redis != nil {
		permissionsJSON, _ := json.Marshal(permissions)
		cache.Redis.Set(ctx, cacheKey, permissionsJSON, config.CacheTTL)
	}

	return permissions, nil
}

// shouldCheckPermissions determines if we should check fresh permissions
func shouldCheckPermissions(c *fiber.Ctx, config *Config) bool {
	method := c.Method()

	switch config.CheckMethod {
	case constant.Grant_Write:
		return method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
	case constant.Grant_Read:
		return method == "GET" || method == "HEAD"
	case constant.Grant_All:
		return true
	default:
		return false
	}
}

// storeClaims stores claims in Fiber context
func storeClaims(c *fiber.Ctx, userID, email string, roles, permissions []string) {
	c.Locals(ClaimsKey, &TokenClaims{
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
	})
	c.Locals(UserIDKey, userID)
	c.Locals(EmailKey, email)
	c.Locals(RolesKey, roles)
	c.Locals(PermissionsKey, permissions)
}

// hasRole checks if user has at least one of the required roles
func hasRole(c *fiber.Ctx, requiredRoles []string) bool {
	userRoles, ok := c.Locals(RolesKey).([]string)
	if !ok {
		return false
	}

	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}

	return false
}

// hasPermission checks if user has at least one of the required permissions
func hasPermission(c *fiber.Ctx, requiredPermissions []string) bool {
	userPermissions, ok := c.Locals(PermissionsKey).([]string)
	if !ok {
		return false
	}

	for _, userPerm := range userPermissions {
		for _, requiredPerm := range requiredPermissions {
			if userPerm == requiredPerm {
				return true
			}
		}
	}

	return false
}

// GetClaims extracts claims from Fiber context
func GetClaims(c *fiber.Ctx) (*TokenClaims, bool) {
	claims, ok := c.Locals(ClaimsKey).(*TokenClaims)
	return claims, ok
}

// GetUserID extracts user ID from Fiber context
func GetUserID(c *fiber.Ctx) (*uuid.UUID, bool) {
	userID, ok := c.Locals(UserIDKey).(string)
	return lib.StrToUUID(userID), ok
}

// GetEmail extracts email from Fiber context
func GetEmail(c *fiber.Ctx) (string, bool) {
	email, ok := c.Locals(EmailKey).(string)
	return email, ok
}

// GetRoles extracts roles from Fiber context
func GetRoles(c *fiber.Ctx) ([]string, bool) {
	roles, ok := c.Locals(RolesKey).([]string)
	return roles, ok
}

// GetPermissions extracts permissions from Fiber context
func GetPermissions(c *fiber.Ctx) ([]string, bool) {
	permissions, ok := c.Locals(PermissionsKey).([]string)
	return permissions, ok
}
