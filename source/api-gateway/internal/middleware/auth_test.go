package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	commonjwt "github.com/eatdetey/letterboxd-replica/source/go-common/pkg/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func TestBearerAuth_RejectsMissingHeader(t *testing.T) {
	app := newAuthTestApp([]byte("secret"))

	req := httptest.NewRequest("GET", "/private", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestBearerAuth_RejectsInvalidOrExpiredToken(t *testing.T) {
	app := newAuthTestApp([]byte("secret"))

	invalidReq := httptest.NewRequest("GET", "/private", nil)
	invalidReq.Header.Set("Authorization", "Bearer not-a-jwt")
	invalidResp, err := app.Test(invalidReq)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	if invalidResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("unexpected status code for invalid token: got %d, want %d", invalidResp.StatusCode, fiber.StatusUnauthorized)
	}

	expiredToken := buildTestToken(t, []byte("secret"), time.Now().Add(-time.Minute))
	expiredReq := httptest.NewRequest("GET", "/private", nil)
	expiredReq.Header.Set("Authorization", "Bearer "+expiredToken)
	expiredResp, err := app.Test(expiredReq)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	if expiredResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("unexpected status code for expired token: got %d, want %d", expiredResp.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestBearerAuth_AllowsValidToken(t *testing.T) {
	secret := []byte("secret")
	app := newAuthTestApp(secret)

	token := buildTestToken(t, secret, time.Now().Add(time.Minute))

	req := httptest.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}

func newAuthTestApp(secret []byte) *fiber.App {
	app := fiber.New()
	app.Use(BearerAuth(secret))
	app.Get("/private", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	return app
}

func buildTestToken(t *testing.T, secret []byte, expiresAt time.Time) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, commonjwt.Claims{
		Id:       1,
		Username: "test",
		Email:    "test@example.com",
		Status:   "active",
		Roles:    []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenStr
}
