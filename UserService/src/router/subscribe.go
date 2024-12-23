package router

import (
	"UserService/src/controller"
	"UserService/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewSubscribeRoute(router fiber.Router, controller *controller.SubscriptionController) {
	app := router.Group("/subscribe", middleware.AuthenticationUser)
	app.Post("/add", controller.AddSubscription)
}