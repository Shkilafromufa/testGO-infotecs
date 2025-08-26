package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"testCaseGO/internal/handler"
	"testCaseGO/internal/service"
	"testCaseGO/internal/setup"
)

func main() {
	err := setup.PrepareEnv()
	{
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func() {
		service.CloseDB(service.Db)
		err := recover()
		if err != nil {
			log.Fatal(err)
		}
	}()
	app := fiber.New(fiber.Config{})
	app.Use(logger.New())
	app.Use(recover())
	handler.RegisterRoutes(app)
	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
