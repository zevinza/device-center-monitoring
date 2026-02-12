package entity

import "github.com/gofiber/fiber/v2"

type BaseController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type BaseControllerPaginated interface {
	GetPaginated(c *fiber.Ctx) error
	GetFiltered(c *fiber.Ctx) error
}
