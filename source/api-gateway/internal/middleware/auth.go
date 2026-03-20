package middleware

import (
	"strings"

	commonjwt "github.com/eatdetey/letterboxd-replica/source/go-common/pkg/jwt"
	"github.com/gofiber/fiber/v3"
)

const authorizationHeader = "Authorization"

func BearerAuth(accessSecret []byte) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, err := extractBearerToken(c.Get(authorizationHeader))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		claims, err := commonjwt.ParseToken(token, accessSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		c.Locals("auth_claims", claims)
		return c.Next()
	}
}

func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(strings.TrimSpace(authHeader), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", fiber.ErrUnauthorized
	}

	return strings.TrimSpace(parts[1]), nil
}
