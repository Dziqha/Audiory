package router

import (
	"UserService/src/controller"
	"UserService/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewReviewRoute(r fiber.Router, controller *controller.ReviewController) {
	app := r.Group("/review", middleware.AuthenticationUser)
	app.Post("/add", controller.AddReview)
	app.Get("/find", controller.GetAllReview)
	app.Get("/findId/:id", controller.GetReviewById)
	app.Put("/update/:id", controller.UpdateReview)
	app.Delete("/delete/:id", controller.DeleteReview)
}