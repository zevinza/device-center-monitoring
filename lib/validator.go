package lib

import (
	"regexp"

	"github.com/google/uuid"
)

// IsValidEmail check is valid email in given string
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Define the email regex pattern
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Match the email against the regex pattern
	return re.MatchString(email)
}

// validate that v is not zero value a.k.a "" or 0
func IsValid[T comparable](v ...T) bool {
	if len(v) == 0 {
		return false
	}
	var zero T
	for _, v := range v {
		if v == zero {
			return false
		}
	}
	return true
}

// IsValidSlices is a function to check if all slices are valid (has length > 0)
func IsValidSlices[T any](slices ...[]T) bool {
	for _, slice := range slices {
		if len(slice) == 0 {
			return false
		}
	}
	return true
}

func IsValidUUID(u uuid.UUID) bool {
	return u != uuid.Nil
}
