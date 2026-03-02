package middleware

import (
	"strings"

	"refina-web-bff/config/env"
	"refina-web-bff/internal/types/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:     false,
				StatusCode: 401,
				Message:    "Authorization token is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:     false,
				StatusCode: 401,
				Message:    "Invalid authorization format",
			})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(env.Cfg.Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:     false,
				StatusCode: 401,
				Message:    "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Status:     false,
				StatusCode: 401,
				Message:    "Invalid token claims",
			})
		}

		userData := dto.UserData{
			ID:       getClaimString(claims, "id"),
			Username: getClaimString(claims, "username"),
			Email:    getClaimString(claims, "email"),
		}

		c.Locals("user_data", userData)
		return c.Next()
	}
}

func getClaimString(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
