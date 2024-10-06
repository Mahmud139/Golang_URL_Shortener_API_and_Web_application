package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *application) routes(r *fiber.App) {
	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	r.Get("/", app.home)
	r.Get("/:url", app.resolveURL)
	r.Post("/shorten-url", limiter.New(limiter.Config{
		Max:                3,
		Expiration:         1 * time.Minute,
		SkipFailedRequests: true,
	}), app.shortenURL)
	r.Get("/app/healthcheck", app.healthCheckHandler)

	r.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
