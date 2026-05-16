package api

import (
	"github.com/gofiber/fiber/v2"

	"chemistry-coach/internal/ai"
	"chemistry-coach/internal/models"
	"chemistry-coach/internal/scenario"
)

var aiClient = ai.NewClient()

func ChatHandler(c *fiber.Ctx) error {

	var req models.ChatRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	prompt := scenario.FirstDateScenario.SystemPrompt + "\n\nUser: " + req.Message

	raw, err := aiClient.Generate(prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.ChatResponse{
		Reply: raw,
	})
}
