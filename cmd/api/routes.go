package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) routes(r *fiber.App) {
	r.Get("", func(c *fiber.Ctx) error {
		res, _ := app.rdb.Get(app.config.db.ctx, "name").Result()
		return c.SendString("Bismillah "+ res)
	})
}