package middleware

import (
	"os"
	"rest-api/models"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authenticated() fiber.Handler {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		panic("No JWT_SECRET environment variable found")
	}
	return jwtware.New(jwtware.Config{
		ContextKey:   "jwt",
		SigningKey:   jwtware.SigningKey{Key: []byte(key)},
		ErrorHandler: jwtError,
	})
}

func AuthUserContext(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		token := c.Locals("jwt").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		err := db.Where("email = ?", email).First(&user).Error
		if err != nil {
			return c.Status(404).SendString("user not found")
		}
		c.Locals("user", user)
		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
