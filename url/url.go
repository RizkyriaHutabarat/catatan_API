package url

import (
	"github.com/RizkyriaHutabarat/catatan_API/controller"
	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	page.Get("/", controller.Home)
	page.Get("/catatan", controller.GetAll)
	page.Get("/catatan/:id", controller.GetCatatanID)
	page.Post("/insert", controller.InsertData)
	page.Put("/update/:id", controller.UpdateData)
	page.Delete("/delete/:id", controller.DeleteCatatan)
	page.Post("/login", controller.Login)
}