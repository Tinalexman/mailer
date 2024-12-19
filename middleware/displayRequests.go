package middleware

import (
	"encoding/json"
	"mailer/logger"

	"github.com/gofiber/fiber/v2"
)

var log = logger.Logger()

func DisplayRequest(c *fiber.Ctx) error {
	log.Println("")
	log.Println("Incoming Request:============================")
	parsedBody := new(interface{})
	c.BodyParser(parsedBody)
	jsonedBody, _ := json.MarshalIndent(*parsedBody, "", "  ")
	IP := c.IP()
	log.Printf("URL: %v", c.Context().URI())
	log.Printf("Method: %v", c.Method())
	log.Printf("Requester's IP: %v", IP)
	log.Printf("Query Params: %v", c.Queries())
	log.Printf("URL Params: %v", c.AllParams())
	if c.Method() == "POST" || c.Method() == "PUT" {
		log.Printf("Body: %s", jsonedBody)
	}
	// if strings.Contains(c.OriginalURL(), "/report/") {
	log.Printf("Headers: %v", c.GetReqHeaders())
	// }
	log.Println("Incoming Request:============================")
	return c.Next()
}
