package routes

import (
	"context"
	"fmt"
	myLogger "mailer/logger"
	"mailer/validation"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"

	"github.com/mailersend/mailersend-go"
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
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		log.Errorln("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func sendEmailWithTarget(data validation.EmailDataWithTarget) error {

	envKey := ""
	if strings.HasPrefix(data.Target, "raso") {
		envKey = "RASO_MAILERSEND_API"
	} else {
		return fmt.Errorf("invalid target: %s", data.Target)
	}

	ms := mailersend.NewMailersend(os.Getenv(envKey))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subject := data.Subject
	text := data.Body

	from := mailersend.From{
		Name:  data.From,
		Email: data.From,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  data.To,
			Email: data.To,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetText(text)

	res, _ := ms.Email.Send(ctx, message)

	fmt.Printf(res.Header.Get("X-Message-Id"))

	log.Println("Email sent successfully")
	return nil
}
