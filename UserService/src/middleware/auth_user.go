package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthenticationUser(c *fiber.Ctx) error {
	tokenSecret := os.Getenv("TOKEN_SECRET_USER")
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	StringToken := strings.TrimPrefix(authHeader, "Bearer ")
	StringToken = strings.TrimSpace(StringToken)
	
	if StringToken == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized header format not supported",
		})
	}


	token, err := jwt.Parse(StringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fiber.ErrUnauthorized
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token: " + err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}
	
	c.Locals("userId", claims["user_id"])
	return c.Next()
}