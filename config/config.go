package config

import "github.com/gofiber/fiber/v2"

var Iteung = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "RIZKYRIA",
	AppName:       "Message Router",
}