package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		reqID := c.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Set("X-Request-Id", reqID)
		c.Locals(requestIDKey, reqID)

		return c.Next()
	}
}
