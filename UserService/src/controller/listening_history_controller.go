package controller

import (
	"UserService/configs"
	"UserService/proto"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type ListeningHistoryController struct{}

func NewListeningHistoryController() *ListeningHistoryController {
	return &ListeningHistoryController{}
}

func (l *ListeningHistoryController) AddListeningHistory(c *fiber.Ctx) error {
	var req models.ListeningHistory

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User
	if err := configs.Database().First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
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

	songClient := proto.NewSongServiceClient(conn)
	log.Println("Song client initialized")

	if req.SongID == 0 {
		log.Println("Song ID is missing in the request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Song ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	songResponse, err := songClient.GetSongById(ctx, &proto.SongRequest{Id: int32(req.SongID)})
	if err != nil {
		log.Printf("Error calling gRPC service: %v", err) 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get song from gRPC service",
			"error":   err.Error(), 
		})
	}

	listeningHistorty :=  models.ListeningHistory{
		UserID: user.ID,
		SongID: int(songResponse.Id),
		DurationPlayed: int(songResponse.Duration),
		PlayedAt: time.Now(),
	}



	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&listeningHistorty); result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response := res.ListeningHistoryRes{
		ID:             req.ID,
		PlayedAt:       req.PlayedAt,
		DurationPlayed: req.DurationPlayed,
		UserID:         req.UserID,
		SongID:         req.SongID,
	}

	return c.JSON(fiber.Map{
		"message": "Listening history added successfully",
		"data":    response,
	})
}

func (l *ListeningHistoryController) GetListeningHistory(c *fiber.Ctx) error {
	starttime := time.Now()
	userID := c.Params("user_id")
	cachekey := os.Getenv("CACHE_KEY_LISTENING_HISTORY_PREFIX") + userID

	val, err := configs.Initialize().Get(context.Background(), cachekey).Result()
	if err  == nil && val != "" {
		var data []res.ListeningHistoryRes
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		c.Set("X-Source", "cache")
		log.Printf("FindListeningHistory from cache: %s", time.Since(starttime).String())
		return c.JSON(fiber.Map{
			"message": "Listening history retrieved successfully",
			"data":    data,
		})
	}

	var histories []models.ListeningHistory
	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("user_id = ?", userID).Find(&histories); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	response := make([]res.ListeningHistoryRes, len(histories))
	for i, history := range histories {
		response[i] = res.ListeningHistoryRes{
			ID:             history.ID,
			PlayedAt:       history.PlayedAt,
			DurationPlayed: history.DurationPlayed,
			UserID:         history.UserID,
			SongID:         history.SongID,
		}
	}
	if len(response) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Listening history not found",
		})
	}

	cachedata, err := json.Marshal(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err = configs.Initialize().Set(context.Background(), cachekey, cachedata, 5*time.Minute).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	c.Set("X-Source", "database")
	log.Printf("FindListeningHistory from database: %s", time.Since(starttime).String())
	return c.JSON(fiber.Map{
		"message": "Listening history retrieved successfully",
		"data":    response,
	})
}
