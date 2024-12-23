package router

import (
	"UserService/src/controller"
	"UserService/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewFavoriteSongRoute(router fiber.Router, controller *controller.FavoriteSongController) {
	app := router.Group("/favorite", middleware.AuthenticationUser)
	app.Post("/add", controller.AddFavoriteSong)
	app.Get("/find", controller.GetAllFavoriteSong)
	app.Get("/findId/:id", controller.GetfavoriteSongById)
	app.Put("/update/:id", controller.UpdateFavoriteSong)
	app.Delete("/delete/:id", controller.DeleteFavoriteSong)
}