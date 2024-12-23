package router

import (
	"UserService/src/controller"

	"github.com/gofiber/fiber/v2"
)


func NewUserRoute(router fiber.Router, controller *controller.UserController) {
	app := router.Group("/user")
	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)
}
