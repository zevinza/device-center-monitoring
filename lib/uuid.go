package lib

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GenUUID() *uuid.UUID {
	u, _ := uuid.NewRandom()
	return &u
}

func StringUUID(u *uuid.UUID) string {
	if u == nil {
		return uuid.Nil.String()
	}
	return u.String()
}

func CompareUUID(u, x *uuid.UUID) bool {
	if u == nil && x == nil {
		return true
	}

	if u == nil || x == nil {
		return false
	}

	return u.String() == x.String()
}

// ParamsUUID return uuid.UUID from params (default value is "id")
func ParamsUUID(c *fiber.Ctx, param ...string) *uuid.UUID {
	if len(param) == 0 {
		param = []string{"id"}
	}
	id, err := uuid.Parse(c.Params(param[0]))
	if nil != err {
		return nil
	}

	return &id
}

// QueryUUID return uuid.UUID from query (default value is "id")
func QueryUUID(c *fiber.Ctx, query ...string) (u []*uuid.UUID) {
	if len(query) == 0 {
		query = []string{"id"}
	}
	for _, q := range strings.Split(c.Query(query[0]), ",") {
		if id, err := uuid.Parse(q); err == nil {
			u = append(u, &id)
		}
	}
	return
}
