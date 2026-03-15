package validation

import (
	"fmt"
	"mailer/logger"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EmailData struct {
	To      string `json:"to" validate:"required"`
	From    string `json:"from" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Config  Config `json:"config" validate:"required"`
}

type EmailDataWithTarget struct {
	To      string `json:"to" validate:"required"`
	From    string `json:"from" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Target  string `json:"target" validate:"required"`
}

type Config struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DefaultHttpResponse struct {
	Status  any    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type IErrorST struct {
	Field        string
	Tag          string
	ErrorMessage string
}

var log = logger.Logger()
var ValidatorObj = validator.New()

func ValidateSendEmail(c *fiber.Ctx) error {
	var errors []*IErrorST
	body := new(EmailData)
	c.BodyParser(&body)

	err := ValidatorObj.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el IErrorST
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.ErrorMessage = err.Error()
			errors = append(errors, &el)
		}

		log.Errorln("(ValidateSendEmail) Validation Failed")
		log.Errorln(err)
		errorResponse := convertValidationErrorToHTTPError(errors)
		errorResponse.Status = fiber.StatusBadRequest
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return c.Next()
}

func ValidateSendEmailWithTarget(c *fiber.Ctx) error {
	var errors []*IErrorST
	body := new(EmailDataWithTarget)
	c.BodyParser(&body)

	err := ValidatorObj.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el IErrorST
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.ErrorMessage = err.Error()
			errors = append(errors, &el)
		}

		log.Errorln("(ValidateSendEmail) Validation Failed")
		log.Errorln(err)
		errorResponse := convertValidationErrorToHTTPError(errors)
		errorResponse.Status = fiber.StatusBadRequest
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if slices.Contains(AllowedTargets, body.Target) == false {
		log.Errorln("(ValidateSendEmailWithoutConfig) Validation Failed: Invalid target query parameter")
		errorResponse := DefaultHttpResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid target query parameter",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return c.Next()
}

func convertValidationErrorToHTTPError(errors []*IErrorST) (errorResponse DefaultHttpResponse) {
	errorToDisplay := errors[0]

	errorResponse.Message = fmt.Sprintf("(%s) %v", errorToDisplay.Field, errorToDisplay.ErrorMessage)

	return errorResponse
}
