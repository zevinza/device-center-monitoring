package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
)

// RecoverSlog catches panics and logs them with slog (keeps the server alive).
func RecoverSlog() fiber.Handler {
	return fiberrecover.New(fiberrecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			reqID, _ := c.Locals(RequestIDKey).(string)
			slog.Error("panic",
				slog.String("request_id", reqID),
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("panic", fmt.Sprint(e)),
			)
		},
	})
}
