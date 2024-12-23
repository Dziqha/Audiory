package router

import (
	"UserService/src/controller"

	"github.com/gofiber/fiber/v2"
)

func NewAdminRoute(router fiber.Router, controller *controller.AdminController) {
	app := router.Group("/admin")
	app.Post("/add", controller.AddAdmin)
	app.Post("/login", controller.LoginAdmin)
	app.Get("/find", controller.GetAllAdmin)
	app.Get("/findId/:id", controller.GetAdminById)
	app.Put("/update/:id", controller.UpdateAdmin)
	app.Delete("/delete/:id", controller.DeleteAdmin)
}