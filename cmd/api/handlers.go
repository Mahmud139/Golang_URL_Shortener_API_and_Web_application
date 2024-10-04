package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/asaskevich/govalidator"
)

type input struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type output struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_remaining"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func (app *application) shortenURL(c *fiber.Ctx) error {
	body := new(input)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON or invalid JSON body",
		})
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incorrect URL, please provide valid URL",
		})
	}

	exist, err := app.checkCustomShort(body.CustomShort)
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if exist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "custom URL has already available",
		})
	}

	

	return nil
}

func (app *application) checkCustomShort(url string) (bool, error) {
	exist, err := app.rdb.Exists(app.config.db.ctx, url).Result()
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (app *application) resolveURL(c *fiber.Ctx) error {
	return nil
}