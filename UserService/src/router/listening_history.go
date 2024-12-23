package router

import (
	"UserService/src/controller"
	"UserService/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewListeningHistoryRoute(router fiber.Router, controller *controller.ListeningHistoryController) {
	app := router.Group("/listening-history", middleware.AuthenticationUser)
	app.Post("/listen", controller.AddListeningHistory)
	app.Get("/:user_id", controller.GetListeningHistory)
}