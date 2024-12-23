package router

import (
	"UserService/src/controller"
	"github.com/gofiber/fiber/v2"
)

func NewRecomendationRoute(router fiber.Router, controller *controller.RecomendationController) {
	app := router.Group("/recomendation")
	app.Get("/:user_id", controller.AddRecommendation)
}