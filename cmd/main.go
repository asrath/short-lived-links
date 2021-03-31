package main

import (
	"log"

	"github.com/asrath/short-lived-links/pkg/api"
	"github.com/asrath/short-lived-links/pkg/config"
	"github.com/asrath/short-lived-links/pkg/home"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	cfg := config.GetConfig()
	engine := html.New("./web/templates", ".html")
	engine.Reload(true)
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("sllConfig", cfg)
		return c.Next()
	})

	app.Static("/static", "web/static")
	app.Get("/", home.Handler)
	app.Get("/:recoveryKey", home.Handler)

	apiPrefix := app.Group("/api")
	v1 := apiPrefix.Group("/v1")
	v1.Post("/paste", api.CreatePasteHandler)
	v1.Get("/paste/:recoveryKey", api.RetrievePasteHandler)

	log.Fatal(app.Listen(":8080"))
}
