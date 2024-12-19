package routes

import (
	myLogger "mailer/logger"
	"mailer/validation"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

var log = myLogger.Logger()

func Email_Controller(api fiber.Router) {

	api.Post("/", validation.ValidateSendEmail, func(c *fiber.Ctx) error {

		var data validation.EmailData
		c.BodyParser(&data)

		if err := sendEmail(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(validation.DefaultHttpResponse{
				Status:  fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(validation.DefaultHttpResponse{
			Status:  fiber.StatusOK,
			Message: "Email Delivered Successfully",
		})

	})

}

func sendEmail(data validation.EmailData) error {
	m := gomail.NewMessage()
	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/plain", data.Body)

	d := gomail.NewDialer(data.Config.Host, data.Config.Port, data.Config.Username, data.Config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Errorln("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
