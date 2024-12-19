package main

import (
	myLogger "mailer/logger"
	"mailer/middleware"
	"mailer/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var log = myLogger.Logger()

func main() {

	app := fiber.New(fiber.Config{Prefork: true, AppName: "Mailer"})

	app.Use(logger.New(logger.Config{
		Format:     "[${time}]-[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Africa/Lagos",
	}))

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(middleware.DisplayRequest)

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			c.Set("Access-Control-Allow-Origin", "*")
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Pagination-Page, X-Pagination-Limit, User-Token, X-Pagination-Totalpages, X-Pagination-Totalcount")
			c.Set("Access-Control-Expose-Headers", "Content-Type, Authorization, X-Pagination-Page, X-Pagination-Limit, User-Token, X-Pagination-Totalpages, X-Pagination-Totalcount")
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Pagination-Page, X-Pagination-Limit, User-Token, X-Pagination-Totalpages, X-Pagination-Totalcount")
		c.Set("Access-Control-Expose-Headers", "Content-Type, Authorization, X-Pagination-Page, X-Pagination-Limit, User-Token, X-Pagination-Totalpages, X-Pagination-Totalcount")
		c.Set("Access-Control-Max-Age", "3600")
		return c.Next()
	})

	app.Route("/api/email", routes.Email_Controller)
	log.Fatal(app.Listen(":2200"))
}
