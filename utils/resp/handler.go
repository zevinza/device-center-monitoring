package resp

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errs, ok := err.(*Response); ok {
		return c.Status(errs.Code).JSON(errs)
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make([]ErrorMessages, len(errs))
		for i, err := range errs {
			errorMessages[i] = ErrorMessages{
				Field:     err.Field(),
				Path:      err.Namespace(),
				Type:      err.Tag(),
				Value:     err.Value(),
				Validator: err.Tag(),
				Message:   err.Error(),
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Code:          fiber.StatusBadRequest,
			Status:        false,
			Message:       "Bad Request",
			ErrorMessages: errorMessages,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Code:    fiber.StatusInternalServerError,
		Status:  false,
		Message: "Internal Server Error",
	})
}
