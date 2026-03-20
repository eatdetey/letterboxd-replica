package middleware

import (
	"log"
	"strings"

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

		// Debug: log Authorization header
		authHeader := c.Get("Authorization")
		path := c.Path()
		if strings.Contains(path, "playlist") {
			log.Printf("DEBUG: path=%s, auth_header_present=%v, auth_prefix=%s", path, authHeader != "", authHeader[:20])
		}

		return c.Next()
	}
}
