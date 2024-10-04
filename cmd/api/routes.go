package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) routes(r *fiber.App) {
	//need to implement rate limiter middleware
	// need to implement panic recover middleware
	// need to implement not found(404) middleware
	r.Get("/:url", app.resolveURL)
	r.Post("/v1/url", app.shortenURL)
}