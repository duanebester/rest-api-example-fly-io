package main

import (
	"os"
	"rest-api/database"
	"rest-api/middleware"
	"rest-api/models"
	"rest-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	db := database.Init()
	app := fiber.New(fiber.Config{
		AppName: "rest-api",
		Prefork: false,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Ping
		err = sqlDB.Ping()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"message": "OK"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		loginInput := new(models.LoginInput)
		if err := c.BodyParser(loginInput); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		// Get user from database
		user := models.User{}
		err := db.Where("email = ?", loginInput.Identity).First(&user).Error
		if err != nil {
			return c.Status(404).SendString("user not found")
		}

		// Compare password
		if !utils.CheckPasswordHash(loginInput.Password, user.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
		}

		// Create JWT token
		token := jwt.New(jwt.SigningMethodHS256)
		key := os.Getenv("JWT_SECRET")
		if key == "" {
			panic("No JWT_SECRET environment variable found")
		}

		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email
		claims["user_id"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte(key))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"access_token": t})
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		userInput := new(models.NewUserInput)
		if err := c.BodyParser(userInput); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		user := models.NewUser(userInput)
		db.Create(&user)
		return c.JSON(user)
	})

	app.Use(middleware.Authenticated())
	app.Use(middleware.AuthUserContext(db))

	app.Get("/currentuser", func(c *fiber.Ctx) error {
		ctxUser := c.Locals("user")
		return c.JSON(ctxUser)
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		ctxUser := c.Locals("user")
		log.Infof("ctxUser %s", ctxUser)

		var users []models.User
		db.Find(&users)
		return c.JSON(users)
	})

	log.Fatal(app.Listen(":8080"))
}
