package middleware

// This file contains usage examples for the authentication middlewares
// Remove this file in production

/*
Example 1: Basic Authentication (no authorization check)
---------------------------------------------------------
app.Get("/api/users/me",
    middleware.OAuth2Authentication(),
    func(c *fiber.Ctx) error {
        userID, _ := middleware.GetUserID(c)
        return c.JSON(fiber.Map{"user_id": userID})
    },
)

Example 2: Authentication + Role-based Authorization
----------------------------------------------------
app.Get("/api/admin/users",
    middleware.OAuth2Authentication(
        middleware.WithRequiredRoles("admin"),
    ),
    func(c *fiber.Ctx) error {
        // Only admins can access
        return c.JSON(fiber.Map{"message": "Admin endpoint"})
    },
)

Example 3: Authentication + Permission-based Authorization
-----------------------------------------------------------
app.Post("/api/posts",
    middleware.OAuth2Authentication(
        middleware.WithRequiredPermissions("write:posts"),
    ),
    func(c *fiber.Ctx) error {
        // Only users with write:posts permission can access
        return c.JSON(fiber.Map{"message": "Post created"})
    },
)

Example 4: Using Convenience Functions
---------------------------------------
app.Get("/api/admin/dashboard",
    middleware.RequireRole("admin"),
    func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Admin dashboard"})
    },
)

app.Delete("/api/posts/:id",
    middleware.RequirePermission("delete:posts"),
    func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Post deleted"})
    },
)

Example 5: Check Permissions Only for Write Operations
-------------------------------------------------------
app.Post("/api/posts",
    middleware.OAuth2Authentication(
        middleware.WithRequiredPermissions("write:posts"),
        middleware.WithCheckMethod("write"), // Only check for POST/PUT/DELETE
        middleware.WithRedisCache(true),
    ),
    func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Post created"})
    },
)

Example 6: Using in Route Groups
--------------------------------
api := app.Group("/api/v1")
api.Use(middleware.OAuth2Authentication()) // Apply to all routes in group

api.Get("/users/me", func(c *fiber.Ctx) error {
    claims, _ := middleware.GetClaims(c)
    return c.JSON(claims)
})

admin := api.Group("/admin", middleware.RequireRole("admin"))
admin.Get("/users", getUsersHandler)
admin.Delete("/users/:id", deleteUserHandler)

Example 7: Accessing User Information in Handlers
--------------------------------------------------
func getUserProfile(c *fiber.Ctx) error {
    // Get user ID
    userID, ok := middleware.GetUserID(c)
    if !ok {
        return c.Status(500).JSON(fiber.Map{"error": "User ID not found"})
    }

    // Get email
    email, _ := middleware.GetEmail(c)

    // Get roles
    roles, _ := middleware.GetRoles(c)

    // Get permissions
    permissions, _ := middleware.GetPermissions(c)

    // Get full claims
    claims, _ := middleware.GetClaims(c)

    return c.JSON(fiber.Map{
        "user_id":     userID,
        "email":       email,
        "roles":       roles,
        "permissions": permissions,
        "claims":      claims,
    })
}

Example 8: Custom Configuration
-------------------------------
app.Post("/api/sensitive-action",
    middleware.OAuth2Authentication(
        middleware.WithJWTSecret("custom-secret"),
        middleware.WithTokenHeader("X-Auth-Token"),
        middleware.WithRequiredPermissions("sensitive:action"),
        middleware.WithRedisCache(true),
        middleware.WithCacheTTL(15 * time.Minute),
    ),
    func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Action completed"})
    },
)

Example 9: Secret Key Authentication for Public APIs
----------------------------------------------------
// Public APIs that don't require user login - use secret key authentication
// Client must send X-API-Key header with the API_SECRET_KEY value
app.Post("/api/public/register",
    middleware.SecretKeyAuthentication(),
    func(c *fiber.Ctx) error {
        // This endpoint is accessible with API secret key
        return c.JSON(fiber.Map{"message": "Registration successful"})
    },
)

app.Post("/api/public/login",
    middleware.SecretKeyAuthentication(),
    func(c *fiber.Ctx) error {
        // Login endpoint - returns JWT token after successful authentication
        return c.JSON(fiber.Map{"token": "jwt_token_here"})
    },
)

// Using in route groups for public APIs
publicAPI := app.Group("/api/public", middleware.SecretKeyAuthentication())
publicAPI.Post("/register", registerHandler)
publicAPI.Post("/login", loginHandler)
publicAPI.Post("/forgot-password", forgotPasswordHandler)

Example 10: Dual Authentication (Secret Key OR JWT)
---------------------------------------------------
// Endpoints that accept both secret key (for public access) and JWT (for authenticated users)
app.Post("/api/reset-password",
    middleware.DualAuthentication(),
    func(c *fiber.Ctx) error {
        // Check which authentication mode was used
        if middleware.IsSecretKeyAuth(c) {
            // Public API access via secret key
            return c.JSON(fiber.Map{"message": "Password reset via public API"})
        } else if middleware.IsJWTAuth(c) {
            // Authenticated user access via JWT
            userID, _ := middleware.GetUserID(c)
            return c.JSON(fiber.Map{"message": "Password reset for user", "user_id": userID})
        }
        return c.JSON(fiber.Map{"message": "Password reset"})
    },
)

// Dual authentication with custom configuration
app.Post("/api/flexible-endpoint",
    middleware.DualAuthentication(
        middleware.WithAPISecretKey("custom-secret-key"),
        middleware.WithSecretKeyHeader("X-Custom-API-Key"),
        middleware.WithJWTSecret("custom-jwt-secret"),
    ),
    func(c *fiber.Ctx) error {
        authMode, _ := middleware.GetAuthMode(c)
        return c.JSON(fiber.Map{"auth_mode": authMode})
    },
)

Example 11: Complete Authentication Flow
----------------------------------------
// 1. Public endpoints (secret key required)
public := app.Group("/api/v1/public", middleware.SecretKeyAuthentication())
public.Post("/register", registerHandler)  // Requires X-API-Key header
public.Post("/login", loginHandler)        // Requires X-API-Key header, returns JWT

// 2. Protected endpoints (JWT required)
protected := app.Group("/api/v1/protected", middleware.OAuth2Authentication())
protected.Get("/profile", profileHandler)           // Requires Authorization: Bearer <jwt>
protected.Post("/logout", logoutHandler)            // Requires Authorization: Bearer <jwt>
protected.Post("/change-password", changePasswordHandler) // Requires Authorization: Bearer <jwt>

// 3. Flexible endpoints (accepts both)
flexible := app.Group("/api/v1/flexible", middleware.DualAuthentication())
flexible.Post("/reset-password", resetPasswordHandler) // Accepts X-API-Key OR Authorization: Bearer <jwt>

Example 12: Checking Authentication Mode in Handlers
-----------------------------------------------------
func myHandler(c *fiber.Ctx) error {
    // Check which authentication method was used
    authMode, ok := middleware.GetAuthMode(c)
    if !ok {
        return c.Status(401).JSON(fiber.Map{"error": "Not authenticated"})
    }

    switch authMode {
    case middleware.AuthModeSecret:
        // Public API access via secret key
        return c.JSON(fiber.Map{
            "message": "Access via secret key",
            "mode": "public",
        })
    case middleware.AuthModeJWT:
        // Authenticated user access
        userID, _ := middleware.GetUserID(c)
        email, _ := middleware.GetEmail(c)
        return c.JSON(fiber.Map{
            "message": "Access via JWT",
            "mode": "authenticated",
            "user_id": userID,
            "email": email,
        })
    default:
        return c.Status(401).JSON(fiber.Map{"error": "Unknown authentication mode"})
    }
}

Example 13: Authorization Middleware - Permission-based
-----------------------------------------------------
// Authorization middleware should be used AFTER authentication middleware
// It checks if the authenticated user has the required permission(s)

import (
    "api/constant/permissions"
    "api/middleware"
)

// Single permission check
app.Put("/api/users/:id",
    middleware.OAuth2Authentication(),           // First authenticate
    middleware.Authorize(permissions.A1),        // Then check permission
    updateUserHandler,
)

// Multiple permissions (OR logic - user needs at least one)
app.Delete("/api/users/:id",
    middleware.OAuth2Authentication(),
    middleware.Authorize(permissions.A3, permissions.A4), // User needs A3 OR A4
    deleteUserHandler,
)

// Using in route groups
userRoutes := app.Group("/api/users", middleware.OAuth2Authentication())
userRoutes.Get("/", getAllUsersHandler)                    // No permission check
userRoutes.Get("/:id", getByIdHandler)                    // No permission check
userRoutes.Post("/", middleware.Authorize(permissions.A1), createUserHandler)
userRoutes.Put("/:id", middleware.Authorize(permissions.A3), updateUserHandler)
userRoutes.Delete("/:id", middleware.Authorize(permissions.A4), deleteUserHandler)

Example 14: Authorization Middleware - Require ALL Permissions
-------------------------------------------------------------
// AuthorizeAll requires the user to have ALL specified permissions (AND logic)
app.Post("/api/sensitive-action",
    middleware.OAuth2Authentication(),
    middleware.AuthorizeAll(permissions.A1, permissions.A8), // User needs A1 AND A8
    sensitiveActionHandler,
)

Example 15: Authorization Middleware - Role-based
-------------------------------------------------
// AuthorizeRole checks if user has the required role(s)
app.Get("/api/admin/dashboard",
    middleware.OAuth2Authentication(),
    middleware.AuthorizeRole("admin", "super_admin"), // User needs admin OR super_admin role
    adminDashboardHandler,
)

Example 16: Authorization Middleware - Flexible (Role OR Permission)
-------------------------------------------------------------------
// AuthorizeAny allows access if user has either the role OR the permission
app.Post("/api/advanced-action",
    middleware.OAuth2Authentication(),
    middleware.AuthorizeAny(
        []string{"admin"},                    // User needs admin role
        []string{permissions.A1},             // OR A1 permission
    ),
    advancedActionHandler,
)

Example 17: Complete Authorization Flow
---------------------------------------
// 1. Public endpoint (no auth)
app.Post("/api/public/register", registerHandler)

// 2. Authenticated endpoint (no specific permission)
app.Get("/api/users/me",
    middleware.OAuth2Authentication(),
    getCurrentUserHandler,
)

// 3. Authenticated + Permission required
app.Put("/api/users/:id",
    middleware.OAuth2Authentication(),
    middleware.Authorize(permissions.A3), // Requires "write:user:update" permission
    updateUserHandler,
)

// 4. Authenticated + Role required
app.Delete("/api/admin/users/:id",
    middleware.OAuth2Authentication(),
    middleware.AuthorizeRole("admin"),
    deleteUserHandler,
)

// 5. Authenticated + Multiple permissions (OR)
app.Post("/api/users",
    middleware.OAuth2Authentication(),
    middleware.Authorize(permissions.A1, permissions.A2), // A1 OR A2
    createUserHandler,
)

// 6. Authenticated + Multiple permissions (AND)
app.Post("/api/users/approve",
    middleware.OAuth2Authentication(),
    middleware.AuthorizeAll(permissions.A1, permissions.A8), // A1 AND A8
    approveUserHandler,
)
*/
