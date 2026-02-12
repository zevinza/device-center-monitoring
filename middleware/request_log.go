package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

// RequestLog logs one line per request with useful server fields.
// It also ensures an X-Request-ID header exists and is available via c.Locals(RequestIDKey).
func RequestLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("X-Request-ID", reqID)
		c.Locals(RequestIDKey, reqID)

		err := c.Next()

		status := c.Response().StatusCode()
		latency := time.Since(start)

		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
			slog.Int("status", status),
			slog.Int64("latency_ms", latency.Milliseconds()),
		}

		if v := c.Locals(UserIDKey); v != nil {
			if s, ok := v.(string); ok && s != "" {
				attrs = append(attrs, slog.String("user_id", s))
			}
		}
		if v := c.Locals(EmailKey); v != nil {
			if s, ok := v.(string); ok && s != "" {
				attrs = append(attrs, slog.String("email", s))
			}
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			slog.LogAttrs(context.Background(), slog.LevelError, "request", attrs...)
			return err
		}

		if status >= 500 {
			slog.LogAttrs(context.Background(), slog.LevelError, "request", attrs...)
			return nil
		}
		if status >= 400 {
			slog.LogAttrs(context.Background(), slog.LevelWarn, "request", attrs...)
			return nil
		}
		slog.LogAttrs(context.Background(), slog.LevelInfo, "request", attrs...)
		return nil
	}
}
