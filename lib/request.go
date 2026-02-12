package lib

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// VALIDATOR validate request body
var VALIDATOR *validator.Validate = validator.New()

// init Register custom validation function
func init() {
	VALIDATOR.RegisterValidation("phone", validatePhone)
	VALIDATOR.RegisterValidation("email", validateEmail)
	VALIDATOR.RegisterValidation("website", validateWebsite)
	VALIDATOR.RegisterValidation("emptyString", validateEmptyString)
	VALIDATOR.RegisterValidation("noWhiteSpace", validateNoWhiteSpace)
}

// BodyParser with validation
func BodyParser(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); nil != err {
		return err
	}

	return VALIDATOR.Struct(payload)
}

// Custom validation function for phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phoneRegex := regexp.MustCompile(`^\d{10,12}$`)
	return phoneRegex.MatchString(fl.Field().String())
}

// Custom validation function for email format
func validateEmail(fl validator.FieldLevel) bool {
	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return emailRegex.MatchString(fl.Field().String())
}

// Custom validation function for website format
func validateWebsite(fl validator.FieldLevel) bool {
	websiteRegex := regexp.MustCompile(`^(http|https):\/\/[^\s/$.?#].[^\s]*$`)
	return websiteRegex.MatchString(fl.Field().String())
}

// Custom validation function for empty string
func validateEmptyString(fl validator.FieldLevel) bool {
	emptyString := regexp.MustCompile(`^\s*$`)
	return !emptyString.MatchString(fl.Field().String())
}

// Custom validation function for no white space
func validateNoWhiteSpace(fl validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), " ") &&
		!strings.HasPrefix(fl.Field().String(), " ") &&
		!strings.HasSuffix(fl.Field().String(), " ")
}
