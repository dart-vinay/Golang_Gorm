package middleware

import "github.com/gofiber/fiber/v2"

func PreMiddleware(c *fiber.Ctx) error {
	return c.Next()
}
