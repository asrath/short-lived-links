package api

import (
	"github.com/asrath/short-lived-links/pkg/config"
	"github.com/asrath/short-lived-links/pkg/storage/paste"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PasteDTO data for the paste received through the API endpoint
type PasteDTO struct {
	Content string    `json:"content" validate:"required,gt=2"`
	TTL     paste.TTL `json:"ttl" validate:"required,oneof=never 1t 1d 2d"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidatePastePayload(dto PasteDTO) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// CreatePasteHandler handler method for the creation of encrypted pastes
func CreatePasteHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	cfg := c.Locals("sllConfig").(*config.Config)

	dto := new(PasteDTO)

	if err := c.BodyParser(dto); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	errors := ValidatePastePayload(*dto)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
	}

	p := paste.New(dto.TTL, cfg.App.PasteStoragePath)
	err := p.Save(dto.Content)
	if err != nil {
		if e, ok := err.(paste.Error); ok {
			c.Status(e.Code)
		} else {
			c.Status(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":      true,
		"recovery_key": p.RecoveryKey,
	})
}

func RetrievePasteHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	cfg := c.Locals("sllConfig").(*config.Config)
	recoveryKey := c.Params("recoveryKey")

	p := paste.Load(recoveryKey, cfg.App.PasteStoragePath)

	if err := p.Retrieve(); err != nil {
		if e, ok := err.(paste.Error); ok {
			c.Status(e.Code)
		} else {
			c.Status(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"content": p.Content,
	})
}
