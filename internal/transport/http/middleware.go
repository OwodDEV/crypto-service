package http

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(c.Context(), "request_id", requestID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func (s *Server) RequestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID, ok := c.UserContext().Value("request_id").(string)
		if !ok {
			requestID = "unknown"
		}

		err := c.Next()

		slog.Info("processed request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.String("remote_ip", c.IP()),
			slog.Duration("latency", time.Duration(time.Since(start).Milliseconds())),
			slog.String("request_id", requestID),
		)
		return err
	}
}
