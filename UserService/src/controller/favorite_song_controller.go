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
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type FavoriteSongController struct{}

func NewFavoriteSongController() *FavoriteSongController {
	return &FavoriteSongController{}
}

func (f *FavoriteSongController) AddFavoriteSong(c *fiber.Ctx) error {
	var req models.FavoriteSongRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return err
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
	log.Println("song response:", songResponse)

	if songResponse.GenreId == 0 {
		log.Println("genre not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "genre not found",
		})
	}

	favoriteSong := models.FavoriteSong{
		UserID:  user.ID,
		SongID:  int(songResponse.Id),
		GenreID: int(songResponse.GenreId),
		AddedAt: time.Now(),
	}

	if err := configs.Database().Transaction(func(tx *gorm.DB) error {
		return tx.Create(&favoriteSong).Error
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response := res.FavoriteSongRes{
		ID:       favoriteSong.ID,
		UserID:   favoriteSong.UserID,
		SongID:   favoriteSong.SongID,
		GenreID:  favoriteSong.GenreID,
		Added_at: favoriteSong.AddedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}


func ( f *FavoriteSongController) GetAllFavoriteSong( c *fiber.Ctx) error {
	starttime := time.Now()
	cachekey := os.Getenv("CACHE_KEY_FAVORITE_SONG_ALL")

	val, err := configs.Initialize().Get(context.Background(), cachekey).Result()
	if err == nil && val != "" {
		var data []res.FavoriteSongRes
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		c.Set("X-Source", "cache")
		log.Printf("FindAllFavoriteSong from cache: %s", time.Since(starttime).String())
		return c.JSON(fiber.Map{
			"message": "Favorite song retrieved successfully",
			"data":    data,
		})
	}
	var favoriteSong []models.FavoriteSong
	var dataResponse []res.FavoriteSongRes

	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Find(&favoriteSong); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get favorite song")
		}

		return nil
	})

	if len(favoriteSong) == 0 {
		return errors.Wrap(err, "favorite song not found")
	}

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}


	for _, findFavotite := range favoriteSong {
		dataResponse = append(dataResponse, res.FavoriteSongRes{
			ID: findFavotite.ID,
			UserID: findFavotite.UserID,
			SongID: findFavotite.SongID,
			GenreID: findFavotite.GenreID,
			Added_at: findFavotite.AddedAt,
		})
	}

	cachedata, err :=  json.Marshal(dataResponse)
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
	log.Printf("FindAllFavoriteSong from database: %s", time.Since(starttime).String())

	return c.JSON(fiber.Map{
		"message": "success fetch data FavoriteSong",
		"data": dataResponse,
	})
}


func (f *FavoriteSongController) GetfavoriteSongById(c *fiber.Ctx) error {
	starttime := time.Now() 
	var findFavoriteSong models.FavoriteSong
	id := c.Params("id")
	cachekey := os.Getenv("CACHE_KEY_FAVORITE_SONG_PREFIX") + id

	val, err := configs.Initialize().Get(context.Background(), cachekey).Result()
	if err == nil && val != "" {
		var data res.FavoriteSongRes
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		c.Set("X-Source", "cache")
		log.Printf("FindFavoriteSongById from cache: %s", time.Since(starttime).String())
		return c.JSON(fiber.Map{
			"message": "Favorite song retrieved successfully",
			"data":    data,
		})
	}


	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		if data := tx.Where("id = ?", id).First(&findFavoriteSong); data.Error != nil {
			return errors.Wrap(data.Error, "failed to find favorite song by id")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	if findFavoriteSong.ID == 0 {
		return errors.Wrap(err, "favorite song not found")
	}

	response := res.FavoriteSongRes{
		ID:       findFavoriteSong.ID,
		UserID:   findFavoriteSong.UserID,
		SongID:   findFavoriteSong.SongID,
		GenreID:  findFavoriteSong.GenreID,
		Added_at: findFavoriteSong.AddedAt,
	}

	cachedata, err :=  json.Marshal(response)
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
	log.Printf("FindFavoriteSongById from database: %s", time.Since(starttime).String())

	return c.JSON(fiber.Map{
		"message": "success fetch data FavoriteSong",
		"data":    response,
	})
}


func (f *FavoriteSongController) UpdateFavoriteSong(c *fiber.Ctx) error {
	var favoriteSong models.FavoriteSong
	var req models.FavoriteSongRequest
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

	if errData := configs.Database().Where("id = ?", id).First(&favoriteSong); errData.Error != nil {
		return errors.Wrap(errData.Error, "failed to find favorite song by id")
	}


	if req.UserID != 0 {
		user := configs.Database().Where("id = ?", req.UserID).First(&models.User{})
		if user.Error != nil {
			return errors.Wrap(user.Error, "failed to find user by id")
		}

		favoriteSong.UserID = req.UserID
	}

	if req.SongID != 0 {
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get song from gRPC service",
				"error":   err.Error(),
			})
		}

		if songResponse.Id == 0 {
			return errors.New("Song not found or invalid ID")
		}
		favoriteSong.SongID = int(songResponse.Id)
	}


	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if data := tx.Save(&favoriteSong); data.Error != nil {
			return errors.Wrap(data.Error, "failed to update favorite song")
		}

		return nil
		
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	responseData := res.FavoriteSongRes{
		ID:       favoriteSong.ID,
		UserID:   favoriteSong.UserID,
		SongID:   favoriteSong.SongID,
		GenreID:  favoriteSong.GenreID,
		Added_at: favoriteSong.AddedAt,
	}

	return c.JSON(fiber.Map{
		"message": "success update favorite song",
		"data":    responseData,
	})
}


func (f *FavoriteSongController)  DeleteFavoriteSong(c *fiber.Ctx) error {
	var favoriteSong models.FavoriteSong
	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if data := tx.Where("id = ?", id).First(&favoriteSong); data.Error != nil {
			return errors.Wrap(data.Error, "failed to find favorite song by id")
		}

		if result := tx.Delete(&favoriteSong); result.Error != nil {
			return errors.Wrap(result.Error, "failed to delete favorite song")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	return c.JSON(fiber.Map{
		"message": "success delete favorite song",
	})
}
