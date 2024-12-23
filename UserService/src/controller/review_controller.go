package controller

import (
	"UserService/configs"
	"UserService/proto"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type ReviewController struct{}

func NewReviewController() *ReviewController {
	return &ReviewController{}
}

func (r *ReviewController) AddReview(c *fiber.Ctx) error {
	var reqData models.ReviewRequest
	var user models.User

	if err := c.BodyParser(&reqData); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}

	if err := utils.ValidateStruct(&reqData); err != nil {
		return errors.Wrap(err, "failed to validate request body")
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to connect to gRPC server",
		})
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}()
	log.Println("Connected to gRPC server")
	
	albumClient := proto.NewAlbumServiceClient(conn)
	log.Println("album client initialized")

	if reqData.AlbumID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Album ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	albumResponse, err := albumClient.GetAlbumById(ctx, &proto.AlbumRequest{Id: int32(reqData.AlbumID)})
	if err != nil {
		log.Printf("Error calling gRPC service: %v", err) 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get song from gRPC service",
			"error":   err.Error(), 
		})
	}

	data :=  models.Review{
		UserID:    reqData.UserID,
		AlbumID:   int(albumResponse.Id),
		Rating:    reqData.Rating,
		Comment:   reqData.Comment,
		CreatedAt: time.Now(),
	}


	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.First(&user, reqData.UserID); result.Error != nil {
			return errors.Wrap(result.Error, "failed to find user")
		}

		if result := tx.Create(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to create review")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	res := res.ReviewResponse{
		ID:        data.ID,
		UserID:    data.UserID,
		AlbumID:   data.AlbumID,
		Rating:    data.Rating,
		Comment:   data.Comment,
		CreatedAt: data.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"message": "Review created successfully",
		"data":    res,
	})
}

func (r *ReviewController) GetReviewById(c *fiber.Ctx) error {
	var data models.Review
	var resData res.ReviewResponse
	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", id).First(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get review by id")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	resData = res.ReviewResponse{
		ID:        data.ID,
		UserID:    data.UserID,
		AlbumID:   data.AlbumID,
		Rating:    data.Rating,
		Comment:   data.Comment,
		CreatedAt: data.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"message": "Successfully get review by id",
		"data":    resData,
	})
}

func (r *ReviewController) GetAllReview(c *fiber.Ctx) error {
	var data []models.Review
	var resData []res.ReviewResponse

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Find(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get all review")
		}
		return nil
	})

	if len(data) == 0 {	
		return errors.Wrap(err, "review not found")
	}

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	for _, result := range data {
		resData = append(resData, res.ReviewResponse{
			ID:        result.ID,
			UserID:    result.UserID,
			AlbumID:   result.AlbumID,
			Rating:    result.Rating,
			Comment:   result.Comment,
			CreatedAt: result.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully get all review",
		"data":    resData,
	})
}


func (r *ReviewController) UpdateReview(c *fiber.Ctx) error {
	var data models.Review
	var reqData models.ReviewRequest
	var resData res.ReviewResponse
	id := c.Params("id")

	if err := c.BodyParser(&reqData); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}

	if result := configs.Database().Where("id = ?", id).First(&data); result.Error != nil {
		return errors.Wrap(result.Error, "failed to get review by id")
	}

	if reqData.UserID != 0 {
		user := configs.Database().Where("id = ?", reqData.UserID).First(&models.User{})
		if user.Error != nil {
			return errors.Wrap(user.Error, "failed to find user by id")
		}
		data.UserID = reqData.UserID
	}

	if reqData.AlbumID != 0 {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to connect to gRPC server for artist service",
			})
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing gRPC connection: %v", err)
			}
		}()
		log.Println("Connected to gRPC server for artist service")

		albumClient := proto.NewAlbumServiceClient(conn)
		log.Println("Album client initialized")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		albumResponse, err := albumClient.GetAlbumById(ctx, &proto.AlbumRequest{Id: int32(reqData.AlbumID)})
		if err != nil {
			log.Printf("Error calling gRPC service for artist: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get album from gRPC service",
				"error":   err.Error(),
			})
		}
		if albumResponse.Id == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Album not found",
			})
		}
		data.AlbumID = int(albumResponse.Id)
	}

	if reqData.Rating != "" {
		data.Rating = reqData.Rating
	}

	if reqData.Comment != "" {
		data.Comment = reqData.Comment
	}


	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Save(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to update review")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	resData = res.ReviewResponse{
		ID:        data.ID,
		UserID:    data.UserID,
		AlbumID:   data.AlbumID,
		Rating:    data.Rating,
		Comment:   data.Comment,
		CreatedAt: data.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update review",
		"data":    resData,
	})
}

func (r *ReviewController) DeleteReview(c *fiber.Ctx) error {
	var data models.Review
	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", id).First(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get review by id")
		}

		if result := tx.Delete(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to delete review")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	return c.JSON(fiber.Map{
		"message": "Successfully delete review",
	})
}