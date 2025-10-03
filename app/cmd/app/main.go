package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"testCaseGO/internal/handler"
	"testCaseGO/internal/setup"
)

// main — точка входа приложения.
// Настраивает окружение, инициализирует Fiber и запускает сервер.
func main() {
	c, err := setup.PrepareEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer c.DB.Close()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(fiberRecover.New())

	handler.NewHTTP(app, c.Svc)

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
