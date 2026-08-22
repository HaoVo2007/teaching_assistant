package middleware

import (
	"strings"

	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtManager *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return response.Fail(c, fiber.StatusUnauthorized, "missing token", "UNAUTHORIZED")
		}

		claims, err := jwtManager.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			return response.Fail(c, fiber.StatusUnauthorized, "invalid token", "UNAUTHORIZED")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func UserIDFromCtx(c *fiber.Ctx) (string, error) {
	id, ok := c.Locals("user_id").(string)
	if !ok || id == "" {
		return "", user.ErrUnauthorized
	}
	return id, nil
}
