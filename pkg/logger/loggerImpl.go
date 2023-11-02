package logger

import "github.com/gofiber/fiber/v2"

type LoggerImpl interface {
	Print() LoggerImpl
	SaveToStorage()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResponse(c any)
}
