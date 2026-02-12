package resp

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code          int             `json:"code"`
	Status        bool            `json:"status"`
	Message       string          `json:"message"`
	Data          any             `json:"data"`
	Page          *int            `json:"page,omitempty"`
	Count         *int            `json:"count,omitempty"`
	Total         *int64          `json:"total,omitempty"`
	TotalPage     *int            `json:"total_page,omitempty"`
	ErrorCode     *string         `json:"error_code,omitempty"`
	ErrorMessages []ErrorMessages `json:"error_messages,omitempty"`
} // @name Response

type ErrorMessages struct {
	Field     string `json:"name"`
	Path      string `json:"path"`
	Type      string `json:"type"`
	Value     any    `json:"value"`
	Validator string `json:"validator"`
	Message   string `json:"message"`
} // @name ErrorMessages

func (r *Response) Error() string {
	return r.Message
}

func OK(c *fiber.Ctx, data ...any) error {
	response := Response{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Success",
	}
	if len(data) > 0 {
		response.Data = data[0]
	}
	return c.JSON(response)
}

func Page[T any](c *fiber.Ctx, data []T, total int64) error {
	page := 1
	if i, _ := strconv.Atoi(c.Query("page")); i > 0 {
		page = i
	}
	count := len(data)
	totalPage := int(total / int64(count))
	if total%int64(count) > 0 {
		totalPage++
	}
	return c.JSON(Response{
		Code:      fiber.StatusOK,
		Status:    true,
		Message:   "Success",
		Data:      data,
		Page:      &page,
		Count:     &count,
		Total:     &total,
		TotalPage: &totalPage,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.JSON(Response{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "Created",
		Data:    data,
	})
}

func ErrorBadRequest(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Bad Request")
	}
	return &Response{
		Code:    fiber.StatusBadRequest,
		Status:  false,
		Message: message[0],
	}
}

func ErrorUnauthorized(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Unauthorized")
	}
	return &Response{
		Code:    fiber.StatusUnauthorized,
		Status:  false,
		Message: message[0],
	}
}

func ErrorForbidden(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Forbidden")
	}
	return &Response{
		Code:    fiber.StatusForbidden,
		Status:  false,
		Message: message[0],
	}
}

func ErrorNotFound(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Not Found")
	}
	return &Response{
		Code:    fiber.StatusNotFound,
		Status:  false,
		Message: message[0],
	}
}

func ErrorNotAllowed(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Not Allowed")
	}
	return &Response{
		Code:    fiber.StatusMethodNotAllowed,
		Status:  false,
		Message: message[0],
	}
}

func ErrorConflict(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Conflict")
	}
	return &Response{
		Code:    fiber.StatusConflict,
		Status:  false,
		Message: message[0],
	}
}

func ErrorInternal(message ...string) error {
	if len(message) == 0 {
		message = append(message, "Internal Server Error")
	}
	return &Response{
		Code:    fiber.StatusInternalServerError,
		Status:  false,
		Message: message[0],
	}
}
