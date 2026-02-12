package middleware

import (
	"api/lib"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

// Authorize checks if the authenticated user has the required permission(s)
// This middleware should be used AFTER OAuth2Authentication middleware
// It supports checking a single permission or multiple permissions (OR logic)
func Authorize(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated (should have permissions in context)
		userPermissions := lib.GetXPermission(c)
		if len(userPermissions) == 0 {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Authentication required"))
		}

		// If no permissions specified, allow access (just checking authentication)
		if len(requiredPermissions) == 0 {
			return c.Next()
		}

		// Check if user has at least one of the required permissions (OR logic)
		hasPermission := false
		for _, userPerm := range userPermissions {
			for _, requiredPerm := range requiredPermissions {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient permissions"))
		}

		return c.Next()
	}
}

// AuthorizeAll checks if the authenticated user has ALL required permissions (AND logic)
// This middleware should be used AFTER OAuth2Authentication middleware
func AuthorizeAll(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated
		userPermissions, ok := c.Locals(PermissionsKey).([]string)
		if !ok || len(userPermissions) == 0 {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Authentication required"))
		}

		// If no permissions specified, allow access
		if len(requiredPermissions) == 0 {
			return c.Next()
		}

		// Check if user has ALL required permissions (AND logic)
		permissionMap := make(map[string]bool)
		for _, userPerm := range userPermissions {
			permissionMap[userPerm] = true
		}

		for _, requiredPerm := range requiredPermissions {
			if !permissionMap[requiredPerm] {
				return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient permissions: missing required permission"))
			}
		}

		return c.Next()
	}
}

// AuthorizeRole checks if the authenticated user has the required role(s)
// This middleware should be used AFTER OAuth2Authentication middleware
func AuthorizeRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated
		userRoles, ok := c.Locals(RolesKey).([]string)
		if !ok || len(userRoles) == 0 {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Authentication required"))
		}

		// If no roles specified, allow access
		if len(requiredRoles) == 0 {
			return c.Next()
		}

		// Check if user has at least one of the required roles
		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient role permissions"))
		}

		return c.Next()
	}
}

// AuthorizeAny checks if the authenticated user has at least one of:
// - One of the required roles, OR
// - One of the required permissions
// This middleware should be used AFTER OAuth2Authentication middleware
func AuthorizeAny(roles []string, permissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated
		_, hasClaims := c.Locals(ClaimsKey).(*TokenClaims)
		if !hasClaims {
			return resp.ErrorHandler(c, resp.ErrorUnauthorized("Authentication required"))
		}

		// Check roles
		if len(roles) > 0 {
			userRoles, ok := c.Locals(RolesKey).([]string)
			if ok {
				for _, userRole := range userRoles {
					for _, requiredRole := range roles {
						if userRole == requiredRole {
							return c.Next()
						}
					}
				}
			}
		}

		// Check permissions
		if len(permissions) > 0 {
			userPermissions, ok := c.Locals(PermissionsKey).([]string)
			if ok {
				for _, userPerm := range userPermissions {
					for _, requiredPerm := range permissions {
						if userPerm == requiredPerm {
							return c.Next()
						}
					}
				}
			}
		}

		return resp.ErrorHandler(c, resp.ErrorForbidden("Insufficient permissions or role"))
	}
}
