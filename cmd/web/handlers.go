package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type input struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

// type output struct {
// 	URL             string        `json:"url"`
// 	ShortenURL      string        `json:"shorten_url"`
// 	Expiry          time.Duration `json:"expiry"`
// 	XRateRemaining  int           `json:"rate_remaining"`
// 	XRateLimitReset time.Duration `json:"rate_limit_reset"`
// }

func (app *application) shortenURL(c *fiber.Ctx) error {
	body := new(input)

	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse form body",
		})
	}

	api_quota, err := app.rdb.Get(app.config.db.ctx, c.IP()).Result()
	if err == redis.Nil {
		err = app.rdb.Set(app.config.db.ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
		if err != nil {
			app.errorLog.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		api_quota = os.Getenv("API_QUOTA")
	} else if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	api_quotaInt, err := strconv.Atoi(api_quota)
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if api_quotaInt <= 0 {
		timeRemaining, err := app.rdb.TTL(app.config.db.ctx, c.IP()).Result()
		if err != nil {
			app.errorLog.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": fmt.Sprintf("rate limit exceeded, reset after %v seconds", timeRemaining),
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

	if !DomainError(body.URL, c) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "host domain can't be shorten",
		})
	}

	body.URL = EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.NewString()[:8]
	} else {
		id = body.CustomShort
	}

	if body.Expiry <= 0 {
		body.Expiry = 24
	}

	err = app.rdb.Set(app.config.db.ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	err = app.rdb.Decr(app.config.db.ctx, c.IP()).Err()
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shortUrl": c.BaseURL() + "/" + id,
	})
}

func (app *application) resolveURL(c *fiber.Ctx) error {
	id := c.Params("url")

	res, err := app.rdb.Get(app.config.db.ctx, id).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "base URL not found against provided shorten URL",
		})
	} else if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Redirect(res, fiber.StatusSeeOther)
}

func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "available",
		"system_info": fiber.Map{
			"environment": app.config.env,
			"version": version,
		},
	})
}

func (app *application) home(c *fiber.Ctx) error {
	tmpl, err := app.Documentation()
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	
	// return tmpl.Execute(c, nil)
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, nil) // Replace 'nil' with data if needed
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	c.Response().Header.Set("content-type", "text/html")

	// Write the buffer content to the response body
	return c.Status(fiber.StatusOK).Send(buf.Bytes())
}
