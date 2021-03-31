package home

import (
	"github.com/asrath/short-lived-links/pkg/config"
	"github.com/asrath/short-lived-links/pkg/storage/paste"
	"github.com/gofiber/fiber/v2"
)

type HomeData struct {
	Title string
}

// Handler for the home page
func Handler(c *fiber.Ctx) error {
	var errMsg string
	var ttl paste.TTL
	cfg := c.Locals("sllConfig").(*config.Config)

	recoveryKey := c.Params("recoveryKey")

	if recoveryKey != "" {
		p := paste.Load(recoveryKey, cfg.App.PasteStoragePath)
		err := p.GetInfo()
		if err != nil {
			errMsg = err.Error()
			return c.Render("error", fiber.Map{
				"title":     cfg.App.Title,
				"logo_text": cfg.App.LogoText,
				"errMsg":    errMsg,
			}, "_base", "layout")
		}
		ttl = p.TTL
	}

	return c.Render("home", fiber.Map{
		"title":       cfg.App.Title,
		"logo_text":   cfg.App.LogoText,
		"recoveryKey": recoveryKey,
		"ttl":         ttl,
	}, "_base", "layout")
}
