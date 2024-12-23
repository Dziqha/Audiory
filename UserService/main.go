package main

import (
	"UserService/configs"
	pb "UserService/proto"
	"UserService/src/controller"
	"UserService/src/middleware"
	"UserService/src/models"
	"UserService/src/router"
	"log"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
    app := fiber.New(fiber.Config{
        IdleTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Prefork:      false,
    })

    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
    configs.Initialize()
    db := configs.Database()
    server := grpc.NewServer()
    pb.RegisterAuthServiceServer(server, &middleware.AuthServiceServer{})
    reflection.Register(server)


    models.MigrateUser(db)
    models.MigrateFavoriteSong(db)
    models.MigrateSubscription(db)
    models.MigrateListeningHistory(db)
    models.MigrateRecomendation(db)
    models.MigrateReview(db)
    models.MigrateAdmin(db)
    models.MigrateRegistration(db)

    usercontroller := controller.NewUserController()
	favoritecontroller := controller.NewFavoriteSongController()
    subscribecontroller := controller.NewSubscriptionController()
    listencontroller := controller.NewListeningHistoryController()
    recomendationcontroller := controller.NewRecomendationController()
    reviewcontroller := controller.NewReviewController()
    admincontroller := controller.NewAdminController()
    registrationcontroller := controller.NewRegistrationController()
    // go func() {
	// 	for {
	// 		if err := subscribecontroller.RemoveExpiredSubscriptions(); err != nil {
	// 			log.Fatalln(err)
	// 		}
    //         configs.Subscribe()
	// 		time.Sleep(24 * time.Hour)
	// 	}
	// }()

    go func() {
        listener, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }
        if err := server.Serve(listener); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    router.NewUserRoute(app, usercontroller)
	router.NewFavoriteSongRoute(app, favoritecontroller)
    router.NewSubscribeRoute(app, subscribecontroller)
    router.NewListeningHistoryRoute(app, listencontroller)
    router.NewRecomendationRoute(app, recomendationcontroller)
    router.NewReviewRoute(app, reviewcontroller)
    router.NewAdminRoute(app, admincontroller)
    router.NewRegistrationRoute(app, registrationcontroller)

    err = app.Listen(":3000")
    if err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}