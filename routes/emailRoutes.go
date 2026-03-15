package routes

import (
	"fmt"
	myLogger "mailer/logger"
	"mailer/validation"
	"os"

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

	api.Post("/target", validation.ValidateSendEmailWithTarget, func(c *fiber.Ctx) error {

		var data validation.EmailDataWithTarget
		c.BodyParser(&data)

		if err := sendEmailWithTarget(data); err != nil {
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

func sendEmailWithTarget(data validation.EmailDataWithTarget) error {
	m := gomail.NewMessage()
	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/plain", data.Body)

	config, err := getConfig(data.Target)
	if err != nil {
		log.Errorln("Error getting config:", err)
		return err
	}

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Errorln("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func getConfig(target string) (validation.Config, error) {
	var config validation.Config

	switch target {
	case "raso-contact":
		config = validation.Config{
			Host:     "rasogroup.co.uk",
			Port:     465,
			Username: os.Getenv("RASO_CONTACT_USERNAME"),
			Password: os.Getenv("RASO_CONTACT_PASSWORD"),
		}
	case "raso-bookings":
		config = validation.Config{
			Host:     "rasogroup.co.uk",
			Port:     465,
			Username: os.Getenv("RASO_BOOKINGS_USERNAME"),
			Password: os.Getenv("RASO_BOOKINGS_PASSWORD"),
		}
	case "raso-services":
		config = validation.Config{
			Host:     "rasogroup.co.uk",
			Port:     465,
			Username: os.Getenv("RASO_SERVICES_USERNAME"),
			Password: os.Getenv("RASO_SERVICES_PASSWORD"),
		}
	case "apata-iye-contact":
		config = validation.Config{
			Host:     "wghp2.wghservers.com",
			Port:     465,
			Username: os.Getenv("APATA_IYE_CONTACT_USERNAME"),
			Password: os.Getenv("APATA_IYE_CONTACT_PASSWORD"),
		}
	default:
		return validation.Config{}, fmt.Errorf("invalid target: %s", target)
	}

	return config, nil
}
