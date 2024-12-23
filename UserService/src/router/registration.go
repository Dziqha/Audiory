package router

import (
	"UserService/src/controller"
	"UserService/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewRegistrationRoute(router fiber.Router, controller *controller.RegistrationController) {
	app := router.Group("/registration", middleware.AuthenticationAdmin)
	app.Post("/add", controller.AddRegister)
	app.Get("/find", controller.GetAllRegistration)
	app.Get("/findId/:id", controller.GetRegistrationById)
	app.Put("/update/:id", controller.UpdateRegistration)
	app.Delete("/delete/:id", controller.DeleteRegistration)
}