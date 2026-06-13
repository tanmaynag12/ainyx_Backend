package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tanmaynag12/ainyx_Backend/internal/logger"
	"go.uber.org/zap"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := uuid.New().String()
		c.Set("X-Request-ID", id)
		c.Locals("request_id", id)
		return c.Next()
	}
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("request_id", c.Locals("request_id").(string)),
		)

		return err
	}
}